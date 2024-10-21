# yt-archive
A mobile-friendly, self-hosted video archiver

# Why Another Video Archiver?
1. **For fun!**
2. Existing solutions didn't work well on mobile devices, such as iPads, due to slow local network speeds.
3. Many required high computing power for advanced features like search and transcoding, which I didn't need.
4. I wanted a lightweight, minimal solution that could run on low-power devices and support adaptive video streaming.

# Key Achievements
- [X] **Runs on Single Board Computers (SBCs)** like Raspberry Pi, Odroid, etc.
- [X] **Mobile-Friendly**: Optimized for playback on mobile devices
- [X] **Readable and Manageable Files**: The archive is organized in a user-friendly format.
- [X] **Simple User Interface**: Easy to use (still deciding on further enhancements).
- [ ] (WIP) Command-line management interface.

# Installation
## Prerequesities
To build the program, ensure you have the following:
* `make`
* `go`
* `node`

You'll also need `FFmpeg` and `yt-dlp` in your system's `PATH` to run the program.

If you're using the Nix package manager with flakes, you can set up the environment with:
```
nix develop
```
## Standalone installation
To build the project locally, run:
```
make all
```
## Docker installation
To build the project using Docker:
```
docker build --tag yt-archive .
```

# Directory Structure
The following directories are used by the program:
* `videos/`: Stores WebM videos.
* `thumbnails/`: Stores video thumbnails and channel avatar images.
* `database/`: Contains metadata (channels, videos, and playlists) in SQLite format.