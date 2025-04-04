# Compiler and flags
CC = gcc
CFLAGS = -Wall -g -std=c99 # -g enables debugging symbols
LDFLAGS =

# Libraries (adjust paths if necessary)
LIBS = -lzmq -lmsgpack-c -lpthread -lrt

# Include paths (adjust if installed in non-standard location)
# Example: CFLAGS += -I/usr/local/include
# Example: LDFLAGS += -L/usr/local/lib

# Target executable name
TARGET = sender

# Source file
SRC = sender.c
SRC += ipc_models.c


# Default target
all: $(TARGET)

# Link the program
$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(SRC) -o build/$(TARGET) $(LDFLAGS) $(LIBS)
	@echo "Build complete: $(TARGET)"

# Clean up build files
clean:
	rm -f $(TARGET)
	@echo "Cleaned up."

.PHONY: all clean