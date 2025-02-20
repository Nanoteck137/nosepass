{ self }:
{ config, lib, pkgs, ... }:
with lib; let
  cfg = config.services.nosepass-web;
in
{
  options.services.nosepass-web = {
    enable = mkEnableOption "Enable the nosepass-web service";

    port = mkOption {
      type = types.port;
      default = 7551;
      description = "port to listen on";
    };

    host = mkOption {
      type = types.str;
      default = "";
      description = "hostname or address to listen on";
    };

    apiAddress = mkOption {
      type = types.str;
      description = "address to the api server";
    };

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.system}.frontend;
      description = "package to use for this service (defaults to the one in the flake)";
    };

    user = mkOption {
      type = types.str;
      default = "nosepass-web";
      description = lib.mdDoc "user to use for this service";
    };

    group = mkOption {
      type = types.str;
      default = "nosepass-web";
      description = lib.mdDoc "group to use for this service";
    };

    openFirewall = mkOption {
      type = types.bool;
      default = false;
      description = "open the ports in the firewall";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.nosepass-web = {
      description = "Frontend for nosepass";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      environment = {
        PORT = "${toString cfg.port}";
        HOST = "${cfg.host}";
        API_ADDRESS = "${cfg.apiAddress}";
        HOST_HEADER = "x-forwarded-host";
        BODY_SIZE_LIMIT = "Infinity";
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;

        ExecStart = "${cfg.package}/bin/nosepass-web";

        Restart = "on-failure";
        RestartSec = "5s";

        ProtectHome = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectProc = "invisible";
        ProtectSystem = "strict";
        RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
      };
    };

    networking.firewall = lib.mkIf cfg.openFirewall {
      allowedTCPPorts = [ cfg.port ];
    };

    users.users = mkIf (cfg.user == "nosepass-web") {
      nosepass-web = {
        group = cfg.group;
        isSystemUser = true;
      };
    };

    users.groups = mkIf (cfg.group == "nosepass-web") {
      nosepass-web = {};
    };
  };
}
