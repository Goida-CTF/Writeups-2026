const std = @import("std");
const W = @import("weights.zig");

const IO_DIM = W.IO_DIM;
const HIDDEN_DIM = W.HIDDEN_DIM;

inline fn relu(x: f32) f32 {
    return @max(0.0, x);
}

fn forward(input: [IO_DIM]f32) [IO_DIM]f32 {
    var h1: [HIDDEN_DIM]f32 = undefined;
    var h2: [HIDDEN_DIM]f32 = undefined;
    var output: [IO_DIM]f32 = undefined;

    //h1 = ReLU(W1 @ input + b1)
    inline for (0..HIDDEN_DIM) |i| {
        var acc: f32 = W.b1[i];
        inline for (0..IO_DIM) |j| {
            acc += W.W1[i][j] * input[j];
        }
        h1[i] = relu(acc);
    }

    //h2 = ReLU(W2 @ h1 + b2)
    inline for (0..HIDDEN_DIM) |i| {
        var acc: f32 = W.b2[i];
        inline for (0..HIDDEN_DIM) |j| {
            acc += W.W2[i][j] * h1[j];
        }
        h2[i] = relu(acc);
    }

    //output = W3 @ h2 + b3
    inline for (0..IO_DIM) |i| {
        var acc: f32 = W.b3[i];
        inline for (0..HIDDEN_DIM) |j| {
            acc += W.W3[i][j] * h2[j];
        }
        output[i] = acc;
    }

    return output;
}

fn derive_key(output: [IO_DIM]f32, input: [IO_DIM]f32) [IO_DIM]u8 {
    var key: [IO_DIM]u8 = undefined;

    for (0..IO_DIM) |i| {
        const out_byte: u8 = @intFromFloat(@abs(output[i]));
        const in_byte: u8 = @intFromFloat(@abs(input[i]));
        key[i] = out_byte ^ in_byte;
    }

    return key;
}

fn solve(key: [IO_DIM]u8, allocator: std.mem.Allocator) !?[]u8 {
    var flag = try allocator.alloc(u8, W.flag_length);

    for (0..W.flag_length) |i| {
        flag[i] = W.encrypted_flag[i] ^ key[i % IO_DIM];
    }

    if (std.mem.startsWith(u8, flag, "goidactf{") and flag[flag.len - 1] == '}') { //TODO
        return flag;
    }

    allocator.free(flag);
    return null;
}

const banner =
    \\
    \\ Привет, я GoidaGPT - самая умная нейросеть, которая может решить любой крипто цтф таск!
    \\ Просто введи 8 чисел float32 через пробел и я дам тебе флаг! (ну или нет)
    \\
;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    var stdout_buf: [8192]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buf);
    const stdout = &stdout_writer.interface;

    var stdin_buf: [1024]u8 = undefined;
    var stdin_reader = std.fs.File.stdin().reader(&stdin_buf);
    const stdin = &stdin_reader.interface;

    try stdout.writeAll(banner);
    try stdout.flush();
    try stdout.writeAll("\n> ");
    try stdout.flush();

    const bare_line = try stdin.takeDelimiter('\n') orelse return error.NoInput;
    const line = std.mem.trim(u8, bare_line, "\r");

    var input: [IO_DIM]f32 = undefined;
    var iter = std.mem.tokenizeScalar(u8, line, ' ');
    var count: usize = 0;

    while (iter.next()) |token| : (count += 1) {
        if (count >= IO_DIM) {
            try stdout.writeAll("[-] Too many values! Expected 8.\n");
            try stdout.flush();
            return;
        }
        input[count] = try std.fmt.parseFloat(f32, token);
    }

    if (count != IO_DIM) {
        try stdout.print("[-] Not enough values! Expected 8, got {d}.\n", .{count});
        try stdout.flush();
        return;
    }

    try stdout.writeAll("\n[*] Решаю...\n");
    try stdout.flush();

    const output = forward(input);

    try stdout.writeAll("[*] Вычисляю флаг...\n");
    try stdout.flush();

    const key = derive_key(output, input);

    try stdout.flush();

    if (try solve(key, allocator)) |flag| {
        defer allocator.free(flag);

        try stdout.writeAll("=== Внимание: обнаружен флаг ===\n\n");
        try stdout.print("FLAG: {s}\n\n", .{flag});
        try stdout.flush();
    } else {
        try stdout.writeAll("[-] Ваш таск - нерешайка\n");
        try stdout.writeAll("[-] Флага тут нет\n");
        try stdout.flush();
    }
}
