package ru.goida.internal.reflect;

import java.lang.reflect.Method;
import ru.goida.crystalstatsmod.CrystalStatsMod;
import ru.goida.internal.security.IntegrityVerifier;

public class HiddenPayloadA {

    private static final byte[] ENCRYPTED_FRAGMENT = {
        0x34,
        0x3C,
        0x3A,
        0x37,
        0x32,
        0x30,
        0x27,
        0x35,
        0x28,
        0x21,
        0x60,
    };

    private static final byte XOR_KEY = 0x53;
    private static final int PAYLOAD_MAGIC = 0xA1B2C3D4;
    private static final int CHAIN_POSITION = 1;
    private static final String NEXT_STAGE =
        "ru.goida.internal.reflect.HiddenPayloadB";

    public static String processPayload(byte[] inputKey, int posMarker) {
        try {
            if (inputKey == null || inputKey.length < 8) {
                return "[PAYLOAD_A_INVALID_INPUT]";
            }

            byte xorKey = XOR_KEY;
            byte[] decryptedFragment = decryptFragment(xorKey);
            byte[] augmentedKey = createAugmentedKey(
                inputKey,
                decryptedFragment
            );

            return invokeNextStage(augmentedKey, decryptedFragment);
        } catch (Exception e) {
            CrystalStatsMod.LOGGER.error("PayloadA processing failed", e);
            return "[PAYLOAD_A_ERROR]";
        }
    }

    private static byte[] decryptFragment(byte xorKey) {
        byte[] result = new byte[ENCRYPTED_FRAGMENT.length];
        for (int i = 0; i < ENCRYPTED_FRAGMENT.length; i++) {
            result[i] = (byte) (ENCRYPTED_FRAGMENT[i] ^ xorKey);
        }
        return result;
    }

    private static byte[] createAugmentedKey(byte[] inputKey, byte[] fragment) {
        byte[] augmented = new byte[32];

        int copyLen = Math.min(inputKey.length, 16);
        System.arraycopy(inputKey, 0, augmented, 0, copyLen);
        System.arraycopy(fragment, 0, augmented, 16, fragment.length);

        augmented[24] = (byte) CHAIN_POSITION;
        augmented[25] = (byte) ENCRYPTED_FRAGMENT.length;
        augmented[26] = (byte) ((PAYLOAD_MAGIC >> 8) & 0xFF);
        augmented[27] = (byte) (PAYLOAD_MAGIC & 0xFF);

        int stageHash = IntegrityVerifier.computeClassHash(
            HiddenPayloadA.class
        );
        augmented[28] = (byte) ((stageHash >> 24) & 0xFF);
        augmented[29] = (byte) ((stageHash >> 16) & 0xFF);
        augmented[30] = (byte) ((stageHash >> 8) & 0xFF);
        augmented[31] = (byte) (stageHash & 0xFF);

        return augmented;
    }

    private static String invokeNextStage(
        byte[] augmentedKey,
        byte[] ourFragment
    ) throws Exception {
        Class<?> nextClass = Class.forName(NEXT_STAGE);
        Method processMethod = nextClass.getDeclaredMethod(
            "processPayload",
            byte[].class,
            byte[].class
        );
        processMethod.setAccessible(true);
        Object result = processMethod.invoke(null, augmentedKey, ourFragment);
        return (result != null) ? result.toString() : null;
    }
}
