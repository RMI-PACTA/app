# Secrets

This directory contains encrypted secret files for credentials that are useful
for developers to have, but aren't necessary for the app to run. We use a tool
called `sops` to manage these encrypted files. You can read more about how
to use it in
[the Developer Handbook](https://siliconally.getoutline.com/s/d984f195-3e5e-410f-bce8-63676496661f#h-sops),
and check out [the root `.sops.yaml` file](/.sops.yaml) to see how files map to
key material.