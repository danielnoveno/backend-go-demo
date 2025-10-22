# Define the name of your Go application executable
APP_NAME := deeply

# Define the main Go source file
MAIN_GO_FILE := main.go

# Target to build the Go application
build:
	go build -o $(APP_NAME) $(MAIN_GO_FILE)

# Target to run the Go application
run: build
	./$(APP_NAME)

# Target to clean up the built executable
clean:
	rm -f $(APP_NAME)