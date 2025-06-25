#!/bin/bash

# GUI Robot - macOS App Bundle Builder
# This script creates a proper macOS app bundle with entitlements

set -e

APP_NAME="GUI Robot"
BUNDLE_ID="com.guirobot.automation"
VERSION="1.0.0"
BINARY_NAME="gui-robot"

# Directories
BUILD_DIR="build/macos"
APP_DIR="$BUILD_DIR/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

echo "🚀 Building GUI Robot macOS App Bundle..."

# Clean previous build
rm -rf "$APP_DIR"

# Create app bundle structure
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

echo "📁 Created app bundle structure"

# Build the binary
echo "🔨 Building GUI Robot binary..."
go build -o "$MACOS_DIR/$BINARY_NAME" ./cmd/gui-robot

echo "✅ Binary built successfully"

# Copy Info.plist
cp "$BUILD_DIR/Info.plist" "$CONTENTS_DIR/"

echo "📄 Copied Info.plist"

# Create PkgInfo
echo "APPLGRBT" > "$CONTENTS_DIR/PkgInfo"

# Copy any resources (icons, etc.)
# TODO: Add app icon if available
# cp icon.icns "$RESOURCES_DIR/"

echo "📦 App bundle created: $APP_DIR"

# Check if we can sign the app
if command -v codesign >/dev/null 2>&1; then
    echo "🔐 Attempting to sign the app..."
    
    # Try to sign with entitlements
    if codesign --force --options runtime --entitlements "$BUILD_DIR/entitlements.plist" --sign - "$APP_DIR" 2>/dev/null; then
        echo "✅ App signed successfully with entitlements"
    else
        echo "⚠️  Signing with entitlements failed, trying basic signing..."
        if codesign --force --sign - "$APP_DIR" 2>/dev/null; then
            echo "✅ App signed with basic signature"
        else
            echo "❌ Code signing failed - app will run unsigned"
        fi
    fi
else
    echo "⚠️  codesign not available - app will run unsigned"
fi

# Make the binary executable
chmod +x "$MACOS_DIR/$BINARY_NAME"

echo ""
echo "🎉 macOS App Bundle created successfully!"
echo "📍 Location: $APP_DIR"
echo ""
echo "To install and run:"
echo "1. Copy '$APP_NAME.app' to /Applications/"
echo "2. Right-click and select 'Open' (first time only)"
echo "3. Grant permissions when prompted"
echo ""
echo "To run from command line:"
echo "  open '$APP_DIR'"
echo ""
echo "To test permissions:"
echo "  '$MACOS_DIR/$BINARY_NAME' -command get_screen_info"

# Create a simple launcher script
cat > "$BUILD_DIR/launch_gui_robot.sh" << 'EOF'
#!/bin/bash
# GUI Robot Launcher
APP_DIR="$(dirname "$0")/GUI Robot.app"
if [ -d "$APP_DIR" ]; then
    open "$APP_DIR"
else
    echo "GUI Robot.app not found!"
    exit 1
fi
EOF

chmod +x "$BUILD_DIR/launch_gui_robot.sh"

echo "🚀 Created launcher script: $BUILD_DIR/launch_gui_robot.sh"
echo ""
echo "✨ Ready to use! The app bundle should now have proper permissions for:"
echo "   • Screen recording"
echo "   • Accessibility control"
echo "   • Input monitoring"
echo "   • Apple Events automation"
