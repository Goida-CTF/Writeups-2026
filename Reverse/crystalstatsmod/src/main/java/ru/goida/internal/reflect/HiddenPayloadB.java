package ru.goida.internal.reflect;

import java.lang.reflect.Method;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.security.IntegrityVerifier;

public class HiddenPayloadB {

    private static final byte[] ENCRYPTED_FRAGMENT = {
        0x11,
        0x1B,
        0x44,
        0x14,
        0x03,
        0x46,
        0x01,
        0x44,
    };

    private static final int PAYLOAD_MAGIC = 0xB2C3D4E5;
    private static final int CHAIN_POSITION = 2;
    private static final String NEXT_STAGE =
        "ru.goida.internal.reflect.FlagAssembler";
    private static final byte BASE_KEY = 0x77;

    public static String processPayload(byte[] augmentedKey, byte[] fragmentA) {
        try {
            if (augmentedKey == null || augmentedKey.length < 24) {
                return "[PAYLOAD_B_INVALID_KEY]";
            }
            if (fragmentA == null || fragmentA.length < 4) {
                return "[PAYLOAD_B_INVALID_FRAG]";
            }

            byte xorKey = deriveXorKey(augmentedKey, fragmentA);
            byte[] decryptedFragment = decryptFragment(xorKey);
            byte[] combinedFragments = combineFragments(
                fragmentA,
                decryptedFragment
            );
            byte[] finalKey = createFinalKey(augmentedKey, combinedFragments);

            return invokeFlagAssembler(finalKey, combinedFragments);
        } catch (Exception e) {
            CrystalStatsMod.LOGGER.error("PayloadB processing failed", e);
            return "[PAYLOAD_B_ERROR]";
        }
    }

    private static byte deriveXorKey(byte[] augmentedKey, byte[] fragmentA) {
        byte key = BASE_KEY;
        key ^= fragmentA[0];
        key ^= fragmentA[0];
        return key;
    }

    private static byte[] decryptFragment(byte xorKey) {
        byte[] result = new byte[ENCRYPTED_FRAGMENT.length];
        for (int i = 0; i < ENCRYPTED_FRAGMENT.length; i++) {
            result[i] = (byte) (ENCRYPTED_FRAGMENT[i] ^ xorKey);
        }
        return result;
    }

    private static byte[] combineFragments(byte[] fragmentA, byte[] fragmentB) {
        byte[] combined = new byte[fragmentA.length + fragmentB.length];
        System.arraycopy(fragmentA, 0, combined, 0, fragmentA.length);
        System.arraycopy(
            fragmentB,
            0,
            combined,
            fragmentA.length,
            fragmentB.length
        );
        return combined;
    }

    private static byte[] createFinalKey(
        byte[] augmentedKey,
        byte[] combinedFragments
    ) {
        int totalLen = 32 + combinedFragments.length + 8;
        byte[] finalKey = new byte[totalLen];

        int copyLen = Math.min(augmentedKey.length, 32);
        System.arraycopy(augmentedKey, 0, finalKey, 0, copyLen);
        System.arraycopy(
            combinedFragments,
            0,
            finalKey,
            32,
            combinedFragments.length
        );

        int metaOffset = 32 + combinedFragments.length;
        finalKey[metaOffset] = (byte) CHAIN_POSITION;
        finalKey[metaOffset + 1] = (byte) ENCRYPTED_FRAGMENT.length;
        finalKey[metaOffset + 2] = (byte) ((PAYLOAD_MAGIC >> 24) & 0xFF);
        finalKey[metaOffset + 3] = (byte) ((PAYLOAD_MAGIC >> 16) & 0xFF);
        finalKey[metaOffset + 4] = (byte) ((PAYLOAD_MAGIC >> 8) & 0xFF);
        finalKey[metaOffset + 5] = (byte) (PAYLOAD_MAGIC & 0xFF);

        int stageHash = IntegrityVerifier.computeClassHash(
            HiddenPayloadB.class
        );
        finalKey[metaOffset + 6] = (byte) ((stageHash >> 8) & 0xFF);
        finalKey[metaOffset + 7] = (byte) (stageHash & 0xFF);

        return finalKey;
    }

    private static String invokeFlagAssembler(
        byte[] finalKey,
        byte[] combinedFragments
    ) throws Exception {
        Class<?> assemblerClass = Class.forName(NEXT_STAGE);
        Method assembleMethod = assemblerClass.getDeclaredMethod(
            "assembleFlag",
            byte[].class,
            byte[].class
        );
        assembleMethod.setAccessible(true);
        Object result = assembleMethod.invoke(
            null,
            finalKey,
            combinedFragments
        );
        return (result != null) ? result.toString() : null;
    }
}
