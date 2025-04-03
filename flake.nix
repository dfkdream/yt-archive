{
  description = "yt-archive: Youtube archiver";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            #backend
            go
            yt-dlp
            ffmpeg
            sqlitebrowser

            #frontend
            nodejs

            imagemagick
          ];
        };

        formatter = pkgs.nixpkgs-fmt;
      });
}
