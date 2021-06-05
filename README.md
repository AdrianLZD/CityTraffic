# City-Traffic


## Windows Only

- Install Chocolatey
  - Run Powershell as an Administrator
  - Run Get-ExecutionPolicy. If it returns Restricted, then run Set-ExecutionPolicy AllSigned or Set-ExecutionPolicy Bypass -Scope Process.
  - Run the following command: "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))"
  - Type "choco" to make sure it installed correctly
- Install Chocolatey Make Package
  -  Run Powershell as an Administrator
  -  Run "choco install make" command 

## Running the Program


| Command | Function |
| ------ | ------ |
| make deps | Installs dependencies |
| make build | Compiles the program |
| make run | Runs the program |
| make run <# cars> <# semaphores> | Runs the program with the specified parameters |
| make | Installs dependencies, compiles the program and runs the program |