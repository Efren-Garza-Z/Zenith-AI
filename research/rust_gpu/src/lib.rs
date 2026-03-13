use std::slice;
use pollster::block_on;
use wgpu::util::DeviceExt;

// Estructura para pasar dimensiones al Shader (debe coincidir con el .wgsl)
#[repr(C)]
#[derive(Copy, Clone, bytemuck::Pod, bytemuck::Zeroable)]
struct Params {
    h: u32,
    w: u32,
    k_size: u32,
    _padding: u32, // Alineación de 16 bytes requerida por WGSL
}

#[no_mangle]
pub extern "C" fn rust_convolve_gpu(
    input_ptr: *const f32,
    output_ptr: *mut f32,
    h: u32,
    w: u32,
    k_ptr: *const f32,
    k_size: u32,
) -> i32 {
    let input_data = unsafe { slice::from_raw_parts(input_ptr, (h * w) as usize) };
    let kernel_data = unsafe { slice::from_raw_parts(k_ptr, (k_size * k_size) as usize) };
    let output_data = unsafe { slice::from_raw_parts_mut(output_ptr, (h * w) as usize) };

    // Ejecutamos la magia de la GPU
    block_on(run_compute(input_data, output_data, h, w, kernel_data, k_size));

    200 // Código de éxito: Procesado por Radeon
}

async fn run_compute(input: &[f32], output: &mut [f32], h: u32, w: u32, kernel: &[f32], k_size: u32) {
    let instance = wgpu::Instance::default();
    let adapter = instance.request_adapter(&wgpu::RequestAdapterOptions::default()).await.unwrap();
    let (device, queue) = adapter.request_device(&wgpu::DeviceDescriptor::default(), None).await.unwrap();

    // Cargamos el shader .wgsl que creaste antes
    let shader = device.create_shader_module(wgpu::ShaderModuleDescriptor {
        label: None,
        source: wgpu::ShaderSource::Wgsl(include_str!("conv.wgsl").into()),
    });

    // --- CREACIÓN DE BUFFERS (VRAM) ---
    let input_buf = device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: None, contents: bytemuck::cast_slice(input), usage: wgpu::BufferUsages::STORAGE,
    });
    let kernel_buf = device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: None, contents: bytemuck::cast_slice(kernel), usage: wgpu::BufferUsages::STORAGE,
    });
    let params_buf = device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: None,
        contents: bytemuck::cast_slice(&[Params { h, w, k_size, _padding: 0 }]),
        usage: wgpu::BufferUsages::UNIFORM,
    });
    let output_buf = device.create_buffer(&wgpu::BufferDescriptor {
        label: None, size: (output.len() * 4) as u64,
        usage: wgpu::BufferUsages::STORAGE | wgpu::BufferUsages::COPY_SRC,
        mapped_at_creation: false,
    });
    // Buffer "Buzón" para leer el resultado desde la CPU
    let staging_buf = device.create_buffer(&wgpu::BufferDescriptor {
        label: None, size: (output.len() * 4) as u64,
        usage: wgpu::BufferUsages::MAP_READ | wgpu::BufferUsages::COPY_DST,
        mapped_at_creation: false,
    });

    // --- PIPELINE Y BIND GROUPS ---
    let compute_pipeline = device.create_compute_pipeline(&wgpu::ComputePipelineDescriptor {
        label: None, layout: None, module: &shader, entry_point: "main",
    });

    let bind_group = device.create_bind_group(&wgpu::BindGroupDescriptor {
        label: None, layout: &compute_pipeline.get_bind_group_layout(0),
        entries: &[
            wgpu::BindGroupEntry { binding: 0, resource: input_buf.as_entire_binding() },
            wgpu::BindGroupEntry { binding: 1, resource: output_buf.as_entire_binding() },
            wgpu::BindGroupEntry { binding: 2, resource: kernel_buf.as_entire_binding() },
            wgpu::BindGroupEntry { binding: 3, resource: params_buf.as_entire_binding() },
        ],
    });

    // --- EJECUCIÓN ---
    let mut encoder = device.create_command_encoder(&wgpu::CommandEncoderDescriptor { label: None });
    {
        let mut cpass = encoder.begin_compute_pass(&wgpu::ComputePassDescriptor { label: None, timestamp_writes: None });
        cpass.set_pipeline(&compute_pipeline);
        cpass.set_bind_group(0, &bind_group, &[]);
        // Lanzamos grupos de 16x16 hilos (coincide con el @workgroup_size en WGSL)
        cpass.dispatch_workgroups((w + 15) / 16, (h + 15) / 16, 1);
    }
    encoder.copy_buffer_to_buffer(&output_buf, 0, &staging_buf, 0, (output.len() * 4) as u64);
    queue.submit(Some(encoder.finish()));

    // --- LEER RESULTADO ---
    let buffer_slice = staging_buf.slice(..);
    let (sender, receiver) = std::sync::mpsc::channel();
    buffer_slice.map_async(wgpu::MapMode::Read, move |v| sender.send(v).unwrap());
    device.poll(wgpu::Maintain::Wait);
    receiver.recv().unwrap().unwrap();

    let data = buffer_slice.get_mapped_range();
    output.copy_from_slice(bytemuck::cast_slice(&data));
}