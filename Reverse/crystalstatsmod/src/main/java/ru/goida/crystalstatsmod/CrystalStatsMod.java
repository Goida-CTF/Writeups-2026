package ru.goida.crystalstatsmod;

import net.fabricmc.api.ModInitializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import ru.goida.internal.core.DiagnosticsCore;

public class CrystalStatsMod implements ModInitializer {

    public static final String MOD_ID = "crystalstatsmod";
    public static final Logger LOGGER = LoggerFactory.getLogger(MOD_ID);

    private static CrystalStatsMod instance;
    private DiagnosticsCore diagnosticsCore;

    @Override
    public void onInitialize() {
        instance = this;
        LOGGER.info("Crystal Stats Mod initializing...");

        // Initialize the diagnostics core system
        diagnosticsCore = new DiagnosticsCore();
        diagnosticsCore.initialize();

        LOGGER.info("Crystal Stats Mod initialized successfully!");
    }

    public static CrystalStatsMod getInstance() {
        return instance;
    }

    public DiagnosticsCore getDiagnosticsCore() {
        return diagnosticsCore;
    }

    public String getSecurityTokenDisplay() {
        if (diagnosticsCore != null && diagnosticsCore.isReady()) {
            return diagnosticsCore.getTokenDisplay();
        }
        return "INITIALIZING...";
    }
}
