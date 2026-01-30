package ru.goida.internal.security;

import java.lang.management.ManagementFactory;
import java.lang.management.RuntimeMXBean;
import java.util.List;
import java.util.Map;

public class DebuggerDetector {

    private static final String[] DEBUG_INDICATORS = {
        "-agentlib:jdwp",
        "-Xdebug",
        "-Xrunjdwp",
        "-javaagent:",
        "suspend=y",
        "transport=dt_socket",
        "-agentpath:",
    };

    private static final String[] DEBUGGER_THREADS = {
        "JDWP",
        "Debugger",
        "Attach Listener",
        "Signal Dispatcher",
    };

    private static volatile Boolean cachedResult = null;
    private static long lastCheckTime = 0;
    private static final long CACHE_DURATION_MS = 5000;

    public static boolean isDebuggerPresent() {
        long currentTime = System.currentTimeMillis();

        if (
            cachedResult != null &&
            (currentTime - lastCheckTime) < CACHE_DURATION_MS
        ) {
            return cachedResult;
        }

        boolean detected = false;

        try {
            detected = detected || checkJvmArguments();
            detected = detected || checkDebuggerThreads();
            detected = detected || checkSystemProperties();
            detected = detected || checkTimingAnomaly();
        } catch (Exception e) {
            detected = false;
        }

        cachedResult = detected;
        lastCheckTime = currentTime;

        return detected;
    }

    private static boolean checkJvmArguments() {
        try {
            RuntimeMXBean runtimeMXBean = ManagementFactory.getRuntimeMXBean();
            List<String> arguments = runtimeMXBean.getInputArguments();

            for (String arg : arguments) {
                String lowerArg = arg.toLowerCase();
                for (String indicator : DEBUG_INDICATORS) {
                    if (lowerArg.contains(indicator.toLowerCase())) {
                        return true;
                    }
                }
            }
        } catch (Exception e) {
            // Ignore
        }

        return false;
    }

    private static boolean checkDebuggerThreads() {
        try {
            Map<Thread, StackTraceElement[]> allThreads =
                Thread.getAllStackTraces();

            for (Thread thread : allThreads.keySet()) {
                String threadName = thread.getName();
                for (String debuggerThread : DEBUGGER_THREADS) {
                    if (threadName.contains(debuggerThread)) {
                        if (
                            debuggerThread.equals("JDWP") ||
                            debuggerThread.equals("Debugger")
                        ) {
                            return true;
                        }
                    }
                }
            }
        } catch (Exception e) {
            // Ignore
        }

        return false;
    }

    private static boolean checkSystemProperties() {
        try {
            String ideDebug = System.getProperty("idea.debugger.dispatch.addr");
            if (ideDebug != null && !ideDebug.isEmpty()) {
                return true;
            }

            String javaDebug = System.getProperty("java.compiler");
            if (
                javaDebug != null && javaDebug.toLowerCase().contains("debug")
            ) {
                return true;
            }
        } catch (Exception e) {
            // Ignore
        }

        return false;
    }

    private static boolean checkTimingAnomaly() {
        try {
            long startTime = System.nanoTime();

            int sum = 0;
            for (int i = 0; i < 1000; i++) {
                sum += i;
            }

            long endTime = System.nanoTime();
            long duration = endTime - startTime;

            if (duration > 100_000_000L) {
                return true;
            }

            if (sum < 0) {
                return true;
            }
        } catch (Exception e) {
            // Ignore
        }

        return false;
    }

    public static int getSecurityStateHash() {
        int hash = 17;
        hash = 31 * hash + DEBUG_INDICATORS.length;
        hash = 31 * hash + DEBUGGER_THREADS.length;
        hash = 31 * hash + DebuggerDetector.class.getDeclaredMethods().length;
        return hash;
    }
}
