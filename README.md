# Vesselizer
Vesselizer is a lightweight container runtime that is designed to run containers with a secure and isolated environment. It is built on top of the Linux kernel's cgroups and namespaces. Vesselizer is written in Go and is designed to be simple and easy to understand.

## How does it work?
Vesselizer uses the Linux kernel's cgroups and namespaces to create an isolated environment for running containers. It creates a new namespace for each container, which isolates the container's processes and filesystem from the host system. It also uses cgroups to limit the resources that the container can use, such as CPU, memory, and disk space.

## Usage
To use Vesselizer, you need to install it as a background service on your system. The following command needs to be run as root:
```vesselizer daemon```

Once the daemon is running, you can use the `vesselizer create` command to create containers.

You need to create a `buildfile` and `entrypoint` file under your project directory which you want to containerize and then run the following command:
```vesselizer create``` inside your project directory.

This will create a container using the provided buildfile taking an alpine image as base and run the container with the specified entrypoint.

**Note that Vesselizer is still in development and is not yet ready for production use.**


## License
Vesselizer is licensed under the MIT License. See the LICENSE file for more information.