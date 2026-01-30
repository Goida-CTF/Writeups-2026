package ru.goida.internal.core;

import java.lang.management.ManagementFactory;
import java.security.MessageDigest;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.security.DebuggerDetector;
import ru.goida.internal.security.IntegrityVerifier;

public class DiagnosticsCore {

    private final AtomicBoolean initialized = new AtomicBoolean(false);
    private final AtomicBoolean ready = new AtomicBoolean(false);
    private final AtomicReference<String> tokenDisplay = new AtomicReference<>(
        "INITIALIZING"
    );

    private Thread initThread;
    private volatile long systemSeed = 0;
    private volatile int integrityHash = 0;

    public void initialize() {
        if (initialized.getAndSet(true)) {
            return;
        }

        initThread = new Thread(
            this::performInitialization,
            "CrystalStats-Init"
        );
        initThread.setDaemon(true);
        initThread.start();
    }

    private void performInitialization() {
        try {
            Thread.sleep(3000);

            CrystalStatsMod.LOGGER.info(
                "DiagnosticsCore: Starting security initialization..."
            );

            if (DebuggerDetector.isDebuggerPresent()) {
                CrystalStatsMod.LOGGER.warn(
                    "DiagnosticsCore: Debug environment detected"
                );
                tokenDisplay.set("SEC_BLOCK");
                return;
            }

            if (!IntegrityVerifier.verifyIntegrity()) {
                CrystalStatsMod.LOGGER.warn(
                    "DiagnosticsCore: Integrity verification failed"
                );
                tokenDisplay.set("INT_FAIL");
                return;
            }

            systemSeed = generateSystemSeed();
            integrityHash = calculateIntegrityHash();

            String token = generateDisplayToken();
            tokenDisplay.set(token);

            ready.set(true);
            CrystalStatsMod.LOGGER.info(
                "DiagnosticsCore: Initialization complete. Token: " + token
            );
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            CrystalStatsMod.LOGGER.error(
                "DiagnosticsCore: Initialization interrupted"
            );
        } catch (Exception e) {
            CrystalStatsMod.LOGGER.error(
                "DiagnosticsCore: Initialization failed",
                e
            );
            tokenDisplay.set("INIT_ERR");
        }
    }

    private long generateSystemSeed() {
        long seed = 0;

        try {
            int methodCount = DiagnosticsCore.class.getDeclaredMethods().length;
            int fieldCount = DiagnosticsCore.class.getDeclaredFields().length;
            int nameHash = DiagnosticsCore.class.getName().hashCode();

            seed =
                ((long) methodCount << 48) |
                ((long) fieldCount << 32) |
                (nameHash & 0xFFFFFFFFL);

            long uptime = ManagementFactory.getRuntimeMXBean().getUptime();
            seed ^= (uptime / 10000) * 10000;
        } catch (Exception e) {
            seed = 0xDEADBEEFL;
        }

        return seed;
    }

    private int calculateIntegrityHash() {
        int hash = 17;

        try {
            String[] classNames = {
                "ru.goida.internal.reflect.ReflectorChainStart",
                "ru.goida.internal.reflect.HiddenPayloadA",
                "ru.goida.internal.reflect.HiddenPayloadB",
                "ru.goida.internal.reflect.FlagAssembler",
            };

            for (String className : classNames) {
                hash = 31 * hash + className.hashCode();
            }

            hash = 31 * hash + DiagnosticsCore.class.getDeclaredFields().length;
        } catch (Exception e) {
            hash = 0xBADC0DE;
        }

        return hash;
    }

    private String generateDisplayToken() {
        try {
            MessageDigest md = MessageDigest.getInstance("MD5");

            long timestamp = System.currentTimeMillis() / 10000;
            String input = String.valueOf(timestamp) + systemSeed;

            byte[] digest = md.digest(input.getBytes());

            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < 4; i++) {
                sb.append(String.format("%02X", digest[i] & 0xFF));
            }

            return sb.toString();
        } catch (Exception e) {
            return "ERR_TOKEN";
        }
    }

    public boolean isReady() {
        return ready.get();
    }

    public String getTokenDisplay() {
        return tokenDisplay.get();
    }
}
