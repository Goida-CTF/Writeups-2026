package ru.goida.internal.security;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.security.MessageDigest;
import java.util.Arrays;

public class IntegrityVerifier {

    private static final int EXPECTED_CORE_SIGNATURE = 0x4D3C2B1A;

    private static volatile boolean verificationComplete = false;
    private static volatile boolean verificationResult = false;

    public static boolean verifyIntegrity() {
        if (verificationComplete) {
            return verificationResult;
        }

        try {
            boolean result = true;
            result = result && verifyClassStructure();
            result = result && verifyMethodSignatures();
            result = result && verifyFieldLayout();
            result = result && verifyCrossReferences();

            verificationResult = result;
            verificationComplete = true;

            return result;
        } catch (Exception e) {
            verificationResult = false;
            verificationComplete = true;
            return false;
        }
    }

    private static boolean verifyClassStructure() {
        try {
            Class<?>[] criticalClasses = {
                IntegrityVerifier.class,
                DebuggerDetector.class,
            };

            int structureHash = 0;
            for (Class<?> clazz : criticalClasses) {
                structureHash ^= clazz.getDeclaredMethods().length;
                structureHash ^= clazz.getDeclaredFields().length;
                structureHash = Integer.rotateLeft(structureHash, 5);
            }

            return structureHash != 0;
        } catch (Exception e) {
            return false;
        }
    }

    private static boolean verifyMethodSignatures() {
        try {
            Class<?> detectorClass = DebuggerDetector.class;
            Method[] methods = detectorClass.getDeclaredMethods();

            boolean hasIsDebuggerPresent = false;
            for (Method method : methods) {
                if (method.getName().equals("isDebuggerPresent")) {
                    hasIsDebuggerPresent = true;
                    if (method.getReturnType() != boolean.class) {
                        return false;
                    }
                    break;
                }
            }

            return hasIsDebuggerPresent;
        } catch (Exception e) {
            return false;
        }
    }

    private static boolean verifyFieldLayout() {
        try {
            Field[] fields = IntegrityVerifier.class.getDeclaredFields();

            int fieldCount = 0;
            boolean hasExpectedSignature = false;

            for (Field field : fields) {
                fieldCount++;
                if (field.getName().equals("EXPECTED_CORE_SIGNATURE")) {
                    hasExpectedSignature = true;
                }
            }

            return fieldCount >= 3 && hasExpectedSignature;
        } catch (Exception e) {
            return false;
        }
    }

    private static boolean verifyCrossReferences() {
        try {
            ClassLoader thisLoader = IntegrityVerifier.class.getClassLoader();
            ClassLoader detectorLoader =
                DebuggerDetector.class.getClassLoader();
            return thisLoader == detectorLoader;
        } catch (Exception e) {
            return false;
        }
    }

    public static int computeClassHash(Class<?> clazz) {
        if (clazz == null) {
            return 0;
        }

        int hash = 17;

        try {
            hash = 31 * hash + clazz.getName().hashCode();

            Method[] methods = clazz.getDeclaredMethods();
            hash = 31 * hash + methods.length;
            for (Method method : methods) {
                hash = 31 * hash + method.getName().hashCode();
                hash = 31 * hash + method.getParameterCount();
            }

            Field[] fields = clazz.getDeclaredFields();
            hash = 31 * hash + fields.length;
        } catch (Exception e) {
            hash = -1;
        }

        return hash;
    }

    public static int getExpectedCoreSignature() {
        return EXPECTED_CORE_SIGNATURE;
    }
}
