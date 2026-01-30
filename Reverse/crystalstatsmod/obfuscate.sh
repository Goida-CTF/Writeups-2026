#!/bin/bash

set -e

BUILD_DIR="build/libs"
PROGUARD_VERSION="7.4.2"
PROGUARD_JAR="proguard-$PROGUARD_VERSION/lib/proguard.jar"
PROGUARD_URL="https://github.com/Guardsquare/proguard/releases/download/v$PROGUARD_VERSION/proguard-$PROGUARD_VERSION.zip"

MOD_JAR=$(find "$BUILD_DIR" -name "crystalstatsmod-*.jar" ! -name "*-sources.jar" ! -name "*-obfuscated.jar" | head -1)

if [ -z "$MOD_JAR" ]; then
    echo "Error: Could not find mod jar in $BUILD_DIR"
    echo "Please run './gradlew build' first"
    exit 1
fi

echo "Found mod jar: $MOD_JAR"

if [ ! -f "$PROGUARD_JAR" ]; then
    echo "Downloading ProGuard $PROGUARD_VERSION..."
    curl -L -o proguard.zip "$PROGUARD_URL"
    unzip -q proguard.zip
    rm proguard.zip
fi

JAVA_HOME=${JAVA_HOME:-$(dirname $(dirname $(readlink -f $(which java))))}

OUTPUT_JAR="${MOD_JAR%.jar}-obfuscated.jar"

java -jar "$PROGUARD_JAR" \
    -injars "$MOD_JAR" \
    -outjars "$OUTPUT_JAR" \
    -libraryjars "$JAVA_HOME/jmods/java.base.jmod(!**.jar,!module-info.class)" \
    -libraryjars "$JAVA_HOME/jmods/java.desktop.jmod(!**.jar,!module-info.class)" \
    -libraryjars "$JAVA_HOME/jmods/java.management.jmod(!**.jar,!module-info.class)" \
    -libraryjars "$JAVA_HOME/jmods/java.logging.jmod(!**.jar,!module-info.class)" \
    -keep 'public class ru.goida.crystalstatsmod.CrystalStatsMod { public *; }' \
    -keep 'public class ru.goida.crystalstatsmod.client.CrystalStatsModClient { public *; }' \
    -keep 'public class ru.goida.crystalstatsmod.client.CrystalStatsHud { public *; }' \
    -keep 'class * implements net.fabricmc.api.ClientModInitializer { public *; }' \
    -keep 'class * implements net.fabricmc.api.ModInitializer { public *; }' \
    -keep '@org.spongepowered.asm.mixin.Mixin class * { *; }' \
    -keep 'class ru.goida.internal.mixin.** { *; }' \
    -keep,allowobfuscation 'public class ru.goida.internal.reflect.** {*; }' \
    -keepattributes 'Signature,InnerClasses,EnclosingMethod' \
    -keepattributes 'RuntimeVisibleAnnotations,RuntimeInvisibleAnnotations' \
    -keepattributes 'RuntimeVisibleParameterAnnotations,RuntimeInvisibleParameterAnnotations' \
    -keepattributes 'AnnotationDefault' \
    -keepclassmembers 'enum * { public static **[] values(); public static ** valueOf(java.lang.String); }' \
    -dontwarn 'net.minecraft.**' \
    -dontwarn 'net.fabricmc.**' \
    -dontwarn 'org.spongepowered.**' \
    -dontwarn 'com.mojang.**' \
    -dontwarn 'org.slf4j.**' \
    -dontwarn 'org.lwjgl.**' \
    -dontwarn 'com.google.**' \
    -dontwarn 'io.netty.**' \
    -dontoptimize \
    -allowaccessmodification \
    -repackageclasses 'ru.goida.o' \
    -printmapping "$BUILD_DIR/proguard-mapping.txt" \
    -verbose

echo ""
echo "Output: $OUTPUT_JAR"
echo "Mapping: $BUILD_DIR/proguard-mapping.txt"
echo ""
