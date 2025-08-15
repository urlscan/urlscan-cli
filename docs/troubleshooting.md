# Troubleshooting

## Keyring

### Using Keyring on Linux

Keyring support for Linux depends on [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring).

```bash
apt install -y gnome-keyring
```

You need to start a D-Bus session manually if you are using headless Linux. For example:

```bash
# start new D-Bus shell
dbus-run-session -- sh
# unlock keyring
echo '...' | gnome-keyring-daemon --unlock
# now you can use keyring
urlscan key set
```
