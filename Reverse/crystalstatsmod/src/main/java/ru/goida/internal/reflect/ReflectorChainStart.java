package ru.goida.internal.reflect;

import java.lang.reflect.Method;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.core.DiagnosticsCore;
import ru.goida.internal.security.DebuggerDetector;
import ru.goida.internal.security.IntegrityVerifier;

public class ReflectorChainStart {

    private static final int CHAIN_MAGIC = 0x47303144;
    private static final String NEXT_STAGE =
        "ru.goida.internal.reflect.HiddenPayloadA";

    private static volatile boolean chainInitialized = false;
    private static volatile long chainTimestamp = 0;

    public static String initiateChain() {
        try {
            if (!performPreflightChecks()) {
                return "[CHAIN_BLOCKED]";
            }

            chainTimestamp = System.currentTimeMillis();
            chainInitialized = true;

            byte[] initialKey = generateInitialKey();
            String result = invokeNextStage(initialKey);

            return postProcessResult(result);
        } catch (Exception e) {
            CrystalStatsMod.LOGGER.error("Chain initiation failed", e);
            return "[CHAIN_ERROR: " + e.getClass().getSimpleName() + "]";
        } finally {
            chainInitialized = false;
        }
    }

    private static boolean performPreflightChecks() {
        if (DebuggerDetector.isDebuggerPresent()) {
            CrystalStatsMod.LOGGER.warn("Chain blocked: debugger detected");
            return false;
        }

        if (!IntegrityVerifier.verifyIntegrity()) {
            CrystalStatsMod.LOGGER.warn(
                "Chain blocked: integrity check failed"
            );
            return false;
        }

        DiagnosticsCore core =
            CrystalStatsMod.getInstance().getDiagnosticsCore();
        if (core == null || !core.isReady()) {
            CrystalStatsMod.LOGGER.warn("Chain blocked: core not ready");
            return false;
        }

        return true;
    }

    private static byte[] generateInitialKey() {
        byte[] key = new byte[16];

        int methodCount = ReflectorChainStart.class.getDeclaredMethods().length;
        int fieldCount = ReflectorChainStart.class.getDeclaredFields().length;

        key[0] = (byte) methodCount;
        key[1] = (byte) fieldCount;
        key[2] = (byte) ((CHAIN_MAGIC >> 24) & 0xFF);
        key[3] = (byte) ((CHAIN_MAGIC >> 16) & 0xFF);
        key[4] = (byte) ((CHAIN_MAGIC >> 8) & 0xFF);
        key[5] = (byte) (CHAIN_MAGIC & 0xFF);

        int nameHash = ReflectorChainStart.class.getSimpleName().hashCode();
        key[6] = (byte) ((nameHash >> 24) & 0xFF);
        key[7] = (byte) ((nameHash >> 16) & 0xFF);
        key[8] = (byte) ((nameHash >> 8) & 0xFF);
        key[9] = (byte) (nameHash & 0xFF);

        key[10] = (byte) NEXT_STAGE.length();

        int securityHash = DebuggerDetector.getSecurityStateHash();
        key[11] = (byte) ((securityHash >> 8) & 0xFF);
        key[12] = (byte) (securityHash & 0xFF);

        int integritySignature = IntegrityVerifier.getExpectedCoreSignature();
        key[13] = (byte) ((integritySignature >> 16) & 0xFF);
        key[14] = (byte) ((integritySignature >> 8) & 0xFF);
        key[15] = (byte) (integritySignature & 0xFF);

        return key;
    }

    private static String invokeNextStage(byte[] keyData) throws Exception {
        Class<?> nextClass = Class.forName(NEXT_STAGE);
        Method processMethod = nextClass.getDeclaredMethod(
            "processPayload",
            byte[].class,
            int.class
        );
        processMethod.setAccessible(true);

        int positionMarker = calculatePositionMarker();
        Object result = processMethod.invoke(null, keyData, positionMarker);

        return (result != null) ? result.toString() : null;
    }

    private static int calculatePositionMarker() {
        int marker = 0;
        try {
            Method[] methods = ReflectorChainStart.class.getDeclaredMethods();
            for (Method m : methods) {
                marker ^= m.getName().hashCode();
            }
        } catch (Exception e) {
            marker = CHAIN_MAGIC;
        }
        return marker;
    }

    private static String postProcessResult(String rawResult) {
        if (rawResult == null || rawResult.isEmpty()) {
            return "[CHAIN_NULL_RESULT]";
        }

        if (rawResult.startsWith("goida{") || rawResult.startsWith("CTF{")) {
            return rawResult;
        }

        return rawResult;
    }
}
