@group(0) @binding(0) var<storage, read> input: array<f32>;
@group(0) @binding(1) var<storage, read_write> output: array<f32>;
@group(0) @binding(2) var<storage, read> kernel_data: array<f32>;

struct Params {
    h: u32,
    w: u32,
    k_size: u32,
};
@group(0) @binding(3) var<uniform> params: Params;

@compute @workgroup_size(16, 16)
fn main(@builtin(global_invocation_id) id: vec3<u32>) {
    let x = id.x;
    let y = id.y;
    let pad = params.k_size / 2u;

    if (y >= pad && y < params.h - pad && x >= pad && x < params.w - pad) {
        var sum: f32 = 0.0;
        for (var ky: u32 = 0u; ky < params.k_size; ky = ky + 1u) {
            for (var kx: u32 = 0u; kx < params.k_size; kx = kx + 1u) {
                let img_idx = (y + ky - pad) * params.w + (x + kx - pad);
                let k_idx = ky * params.k_size + kx;
                sum = sum + input[img_idx] * kernel_data[k_idx];
            }
        }
        output[y * params.w + x] = sum;
    }
}