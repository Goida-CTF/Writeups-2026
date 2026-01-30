package ru.goida.crystalstatsmod.client;

import net.fabricmc.api.EnvType;
import net.fabricmc.api.Environment;
import net.minecraft.client.MinecraftClient;
import net.minecraft.client.font.TextRenderer;
import net.minecraft.client.gui.DrawContext;
import net.minecraft.client.render.RenderTickCounter;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.core.DiagnosticsCore;

@Environment(EnvType.CLIENT)
public class CrystalStatsHud {

    private boolean visible = true;
    private long lastUpdateTime = 0;
    private String cachedToken = "LOADING...";

    private static final int COLOR_TITLE = 0xFF55FF55;      // Green
    private static final int COLOR_LABEL = 0xFFAAAAAA;      // Gray
    private static final int COLOR_VALUE = 0xFFFFFFFF;      // White
    private static final int COLOR_TOKEN = 0xFFFFFF55;      // Yellow
    private static final int COLOR_READY = 0xFF55FF55;      // Green
    private static final int COLOR_NOT_READY = 0xFFFF5555;  // Red
    private static final int BACKGROUND_COLOR = 0x80000000; // Semi-transparent black

    public void render(DrawContext context, RenderTickCounter tickCounter) {
        if (!visible) return;

        MinecraftClient client = MinecraftClient.getInstance();

        if (client.getDebugHud().shouldShowDebugHud()) return;

        if (client.world == null) return;

        TextRenderer textRenderer = client.textRenderer;
        int x = 5;
        int y = 5;
        int lineHeight = 10;
        int padding = 3;

        int boxWidth = 130;
        int boxHeight = 60;

        context.fill(x - padding, y - padding, x + boxWidth + padding, y + boxHeight + padding, BACKGROUND_COLOR);

        context.drawTextWithShadow(textRenderer, "§a§lCrystal Stats", x, y, COLOR_TITLE);
        y += lineHeight + 2;

        int fps = client.getCurrentFps();
        String fpsColor = fps >= 60 ? "§a" : (fps >= 30 ? "§e" : "§c");
        context.drawTextWithShadow(textRenderer, "§7FPS: " + fpsColor + fps, x, y, COLOR_LABEL);
        y += lineHeight;

        Runtime runtime = Runtime.getRuntime();
        long usedMemory = (runtime.totalMemory() - runtime.freeMemory()) / 1024 / 1024;
        long maxMemory = runtime.maxMemory() / 1024 / 1024;
        int memPercent = (int) ((usedMemory * 100) / maxMemory);
        String memColor = memPercent < 70 ? "§a" : (memPercent < 90 ? "§e" : "§c");
        context.drawTextWithShadow(textRenderer, "§7Mem: " + memColor + usedMemory + "/" + maxMemory + " MB", x, y, COLOR_LABEL);
        y += lineHeight;

        updateToken();
        context.drawTextWithShadow(textRenderer, "§7Token: §e" + cachedToken, x, y, COLOR_LABEL);
        y += lineHeight;

        DiagnosticsCore core = CrystalStatsMod.getInstance() != null ?
            CrystalStatsMod.getInstance().getDiagnosticsCore() : null;
        boolean ready = core != null && core.isReady();
        String statusSymbol = ready ? "§a●" : "§c●";
        String statusText = ready ? " §7Ready" : " §7Init...";
        context.drawTextWithShadow(textRenderer, statusSymbol + statusText, x, y, COLOR_LABEL);
    }

    private void updateToken() {
        long currentTime = System.currentTimeMillis();

        if (currentTime - lastUpdateTime > 500) {
            lastUpdateTime = currentTime;

            if (CrystalStatsMod.getInstance() != null) {
                cachedToken = CrystalStatsMod.getInstance().getSecurityTokenDisplay();
            }
        }
    }

    public void toggleVisibility() {
        visible = !visible;
    }

    public boolean isVisible() {
        return visible;
    }

    public void setVisible(boolean visible) {
        this.visible = visible;
    }
}
