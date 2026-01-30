package ru.goida.crystalstatsmod.client;

import net.fabricmc.api.ClientModInitializer;
import net.fabricmc.api.EnvType;
import net.fabricmc.api.Environment;
import net.fabricmc.fabric.api.client.command.v2.ClientCommandManager;
import net.fabricmc.fabric.api.client.command.v2.ClientCommandRegistrationCallback;
import net.fabricmc.fabric.api.client.rendering.v1.HudRenderCallback;
import net.minecraft.text.Text;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.core.DiagnosticsCore;

@Environment(EnvType.CLIENT)
public class CrystalStatsModClient implements ClientModInitializer {

    private CrystalStatsHud hudRenderer;

    @Override
    public void onInitializeClient() {
        CrystalStatsMod.LOGGER.info("Crystal Stats Client initializing...");

        hudRenderer = new CrystalStatsHud();
        HudRenderCallback.EVENT.register(hudRenderer::render);
        registerCommands();

        CrystalStatsMod.LOGGER.info("Crystal Stats Client initialized!");
    }

    private void registerCommands() {
        ClientCommandRegistrationCallback.EVENT.register(
            (dispatcher, registryAccess) -> {
                dispatcher.register(
                    ClientCommandManager.literal("crystalstats")
                        .then(
                            ClientCommandManager.literal("status").executes(
                                context -> {
                                    DiagnosticsCore core =
                                        CrystalStatsMod.getInstance().getDiagnosticsCore();
                                    boolean ready =
                                        core != null && core.isReady();

                                    context
                                        .getSource()
                                        .sendFeedback(
                                            Text.literal(
                                                "§6[CrystalStats] §7System Status: " +
                                                    (ready
                                                        ? "§aONLINE"
                                                        : "§cOFFLINE")
                                            )
                                        );
                                    context
                                        .getSource()
                                        .sendFeedback(
                                            Text.literal(
                                                "§6[CrystalStats] §7Security Token: §e" +
                                                    CrystalStatsMod.getInstance().getSecurityTokenDisplay()
                                            )
                                        );
                                    return 1;
                                }
                            )
                        )
                        .then(
                            ClientCommandManager.literal("toggle").executes(
                                context -> {
                                    hudRenderer.toggleVisibility();
                                    boolean visible = hudRenderer.isVisible();
                                    context
                                        .getSource()
                                        .sendFeedback(
                                            Text.literal(
                                                "§6[CrystalStats] §7HUD " +
                                                    (visible
                                                        ? "§aenabled"
                                                        : "§cdisabled")
                                            )
                                        );
                                    return 1;
                                }
                            )
                        )
                );
            }
        );
    }
}
