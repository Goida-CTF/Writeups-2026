package ru.goida.internal.reflect;

import ru.goida.crystalstatsmod.CrystalStatsMod;

public class FlagAssembler {

    private static final byte[] ENCRYPTED_FINAL_FRAGMENT = {
        0x41,
        0x73,
        0x2F,
        0x66,
        0x2F,
        0x70,
        0x41,
        0x73,
        0x2A,
        0x6D,
        0x6A,
        0x2D,
        0x6C,
        0x63,
    };

    private static final int ASSEMBLY_MAGIC = 0xF1A6F1A6;
    private static final byte FINAL_XOR_KEY = 0x1E;
    private static final String FLAG_PREFIX = "goida{";
    private static final String FLAG_SUFFIX = "}";

    public static String assembleFlag(
        byte[] finalKey,
        byte[] combinedFragments
    ) {
        try {
            if (finalKey == null || finalKey.length < 32) {
                return "[ASSEMBLER_INVALID_KEY]";
            }
            if (combinedFragments == null || combinedFragments.length < 8) {
                return "[ASSEMBLER_INVALID_FRAGS]";
            }

            byte[] finalFragment = decryptFinalFragment();
            byte[] completeFlag = combineAllFragments(
                combinedFragments,
                finalFragment
            );
            String flagString = new String(completeFlag);

            if (!validateFlagFormat(flagString)) {
                flagString = reconstructFlag(combinedFragments, finalFragment);
            }

            CrystalStatsMod.LOGGER.info("Flag assembly complete");
            return flagString;
        } catch (Exception e) {
            CrystalStatsMod.LOGGER.error("Flag assembly failed", e);
            return "[ASSEMBLER_ERROR: " + e.getMessage() + "]";
        }
    }

    private static byte[] decryptFinalFragment() {
        byte[] result = new byte[ENCRYPTED_FINAL_FRAGMENT.length];
        for (int i = 0; i < ENCRYPTED_FINAL_FRAGMENT.length; i++) {
            result[i] = (byte) (ENCRYPTED_FINAL_FRAGMENT[i] ^ FINAL_XOR_KEY);
        }
        return result;
    }

    private static byte[] combineAllFragments(
        byte[] combinedAB,
        byte[] finalFrag
    ) {
        byte[] complete = new byte[combinedAB.length + finalFrag.length];
        System.arraycopy(combinedAB, 0, complete, 0, combinedAB.length);
        System.arraycopy(
            finalFrag,
            0,
            complete,
            combinedAB.length,
            finalFrag.length
        );
        return complete;
    }

    private static boolean validateFlagFormat(String flag) {
        if (flag == null || flag.isEmpty()) {
            return false;
        }

        boolean hasPrefix = flag.startsWith(FLAG_PREFIX);
        boolean hasSuffix = flag.endsWith(FLAG_SUFFIX);
        boolean validLength = flag.length() >= 10 && flag.length() <= 100;

        boolean allPrintable = true;
        for (char c : flag.toCharArray()) {
            if (c < 0x20 || c > 0x7E) {
                allPrintable = false;
                break;
            }
        }

        return hasPrefix && hasSuffix && validLength && allPrintable;
    }

    private static String reconstructFlag(byte[] combinedAB, byte[] finalFrag) {
        StringBuilder sb = new StringBuilder();

        for (byte b : combinedAB) {
            if (b >= 0x20 && b <= 0x7E) {
                sb.append((char) b);
            }
        }

        for (byte b : finalFrag) {
            if (b >= 0x20 && b <= 0x7E) {
                sb.append((char) b);
            }
        }

        return sb.toString();
    }
}
