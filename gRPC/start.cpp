#include <windows.h>
#include <iostream>

int main() {
    STARTUPINFO si1 = { sizeof(STARTUPINFO) };
    STARTUPINFO si2 = { sizeof(STARTUPINFO) };
    PROCESS_INFORMATION pi1, pi2;

    // Commands (must be writable)
    char cmd1[] = "go run C:\\Users\\User\\Desktop\\Bain\\Go\\journal\\main.go";
    char cmd2[] = "go run C:\\Users\\User\\Desktop\\Bain\\Go\\Auth\\server\\server.go";

    // Start Journal Server (non-blocking)
    if (!CreateProcess(
        NULL,       // No module name (use command line)
        cmd1,      // Command line
        NULL,      // Process handle not inheritable
        NULL,      // Thread handle not inheritable
        FALSE,     // Set handle inheritance to FALSE
        CREATE_NEW_CONSOLE,  // New console window
        NULL,      // Use parent's environment
        NULL,      // Use parent's directory
        &si1,      // Pointer to STARTUPINFO
        &pi1       // Pointer to PROCESS_INFORMATION
    )) {
        std::cerr << "Failed to start Journal Server" << std::endl;
        return 1;
    }

    // Start Auth Server (non-blocking)
    if (!CreateProcess(
        NULL,
        cmd2,
        NULL,
        NULL,
        FALSE,
        CREATE_NEW_CONSOLE,
        NULL,
        NULL,
        &si2,
        &pi2
    )) {
        std::cerr << "Failed to start Auth Server" << std::endl;
        return 1;
    }

    std::cout << "Both servers are running in separate consoles." << std::endl;

    // Optional: Wait for user input before closing manager
    std::cout << "Press Enter to exit...";
    std::cin.ignore();

    // Clean up handles (optional, since processes run independently)
    CloseHandle(pi1.hProcess);
    CloseHandle(pi1.hThread);
    CloseHandle(pi2.hProcess);
    CloseHandle(pi2.hThread);

    return 0;
}
