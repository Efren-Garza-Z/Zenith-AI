const std = @import("std");

export fn zig_convolve_fast(
    input_ptr: [*]const f64,
    output_ptr: [*]f64,
    img_h: i32,
    img_w: i32,
    k_ptr: [*]const f64,
    k_size: i32,
) void {
    const h = @as(usize, @intCast(img_h));
    const w = @as(usize, @intCast(img_w));
    const ks = @as(usize, @intCast(k_size));
    const pad = ks / 2;

    // Convertimos a slices para ayudar al optimizador (bounds checking en Debug, velocidad en Release)
    const input = input_ptr[0 .. h * w];
    const output = output_ptr[0 .. h * w];
    const kernel = k_ptr[0 .. ks * ks];

    var y: usize = pad;
    while (y < h - pad) : (y += 1) {
        var x: usize = pad;
        while (x < w - pad) : (x += 1) {
            var sum: f64 = 0;
            var ky: usize = 0;
            while (ky < ks) : (ky += 1) {
                const img_row_offset = (y + ky - pad) * w;
                const k_row_offset = ky * ks;

                var kx: usize = 0;
                while (kx < ks) : (kx += 1) {
                    sum += input[img_row_offset + (x + kx - pad)] * kernel[k_row_offset + kx];
                }
            }
            output[y * w + x] = sum;
        }
    }
}

export fn zig_gaussian_blur(
    input_ptr: [*]const f64,
    output_ptr: [*]f64,
    img_h: i32,
    img_w: i32,
    k_ptr: [*]const f64,
    k_size: i32,
) void {
    const h = @as(usize, @intCast(img_h));
    const w = @as(usize, @intCast(img_w));
    const ks = @as(usize, @intCast(k_size));
    const pad = ks / 2;

    const input = input_ptr[0 .. h * w];
    const output = output_ptr[0 .. h * w];
    const kernel = k_ptr[0 .. ks * ks];

    var y: usize = pad;
    while (y < h - pad) : (y += 1) {
        var x: usize = pad;
        while (x < w - pad) : (x += 1) {
            var sum: f64 = 0;
            var ky: usize = 0;
            while (ky < ks) : (ky += 1) {
                const img_row_offset = (y + ky - pad) * w;
                const k_row_offset = ky * ks;
                var kx: usize = 0;
                while (kx < ks) : (kx += 1) {
                    sum += input[img_row_offset + (x + kx - pad)] * kernel[k_row_offset + kx];
                }
            }
            output[y * w + x] = sum;
        }
    }
}

export fn zig_threshold(
    input_ptr: [*]const f64,
    output_ptr: [*]f64,
    h: i32,
    w: i32,
    threshold: f64,
) void {
    const size = @as(usize, @intCast(h * w));
    const input = input_ptr[0..size];
    const output = output_ptr[0..size];

    var i: usize = 0;
    while (i < size) : (i += 1) {
        if (input[i] >= threshold) {
            output[i] = 1.0;
        } else {
            output[i] = 0.0;
        }
    }
}