# ProGuard rules for CrystalStatsMod CTF Challenge
# Obfuscates internal classes while keeping entry points readable

# Keep the main mod class readable for initial analysis
-keep public class ru.goida.crystalstatsmod.CrystalStatsMod {
    public *;
}

# Keep Fabric API interfaces
-keep class * implements net.fabricmc.api.ClientModInitializer {
    public *;
}
-keep class * implements net.fabricmc.api.ModInitializer {
    public *;
}
-keep class * implements net.fabricmc.api.DedicatedServerModInitializer {
    public *;
}

# Keep Mixin classes and their structure (required for Fabric)
-keep @org.spongepowered.asm.mixin.Mixin class * {
    *;
}
-keep class ru.goida.internal.mixin.** {
    *;
}

# Keep HUD renderer class name but obfuscate internals
-keep class ru.goida.crystalstatsmod.client.CrystalStatsHud {
    public *;
}

# Obfuscate internal security classes - rename but keep methods working
-keep,allowobfuscation class ru.goida.internal.security.** {
    *;
}

# Obfuscate reflection chain classes
-keep,allowobfuscation class ru.goida.internal.reflect.** {
    *;
}

# Obfuscate core classes
-keep,allowobfuscation class ru.goida.internal.core.** {
    *;
}

# Keep attributes needed for reflection and mixins
-keepattributes Signature,InnerClasses,EnclosingMethod
-keepattributes RuntimeVisibleAnnotations,RuntimeInvisibleAnnotations
-keepattributes RuntimeVisibleParameterAnnotations,RuntimeInvisibleParameterAnnotations
-keepattributes AnnotationDefault
-keepattributes SourceFile,LineNumberTable

# Keep enum classes
-keepclassmembers enum * {
    public static **[] values();
    public static ** valueOf(java.lang.String);
}

# Don't warn about missing Minecraft/Fabric classes
-dontwarn net.minecraft.**
-dontwarn net.fabricmc.**
-dontwarn org.spongepowered.**
-dontwarn com.mojang.**
-dontwarn org.slf4j.**
-dontwarn org.lwjgl.**
-dontwarn com.google.**
-dontwarn io.netty.**

# Don't optimize - we want the challenge solvable
-dontoptimize

# Allow access modification for obfuscation
-allowaccessmodification

# Repackage obfuscated classes
-repackageclasses 'ru.goida.internal.o'

# Use mixed-case class names for slight obfuscation
-classobfuscationdictionary proguard-dictionary.txt
-obfuscationdictionary proguard-dictionary.txt

# Output mapping file for verification
-printmapping build/libs/proguard-mapping.txt
