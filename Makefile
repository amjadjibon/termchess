.PHONY: build run clean

# Directories and file names
BIN_FOLDER := bin
BIN_NAME := termchess
BINARY := $(BIN_FOLDER)/$(BIN_NAME)

# Go build flags for optimization
GO_BUILD_FLAGS := -ldflags="-s -w" # Strip debugging information and reduce binary size
GO_BUILD_FLAGS += -trimpath # Remove file system paths from the compiled binary
GO_BUILD_FLAGS += -o $(BINARY) # Output file

# Environment variables for optimization
BUILD_ENV := CGO_ENABLED=0

# Main Go file
MAIN_FILE := main.go

# Build target
build:
	@echo "Building... @$(BINARY)"
	@mkdir -p $(BIN_FOLDER)
	@$(BUILD_ENV) go build $(GO_BUILD_FLAGS) ${MAIN_FILE}

# Run target
run: build
	@echo "Starting... @$(BINARY)"
	@$(BINARY)
	@echo "Exiting... @$(BINARY)"

# Clean target to remove generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_FOLDER)
