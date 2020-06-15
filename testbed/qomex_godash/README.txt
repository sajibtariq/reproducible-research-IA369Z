Steps for goDASH/goDASHbed install, via install script

1.	make the installation script executable:

chmod +x <script.sh>

2.	run script:

./<script.sh>

Note:
a.	for some installs, we get a "E: Could not get lock /var/lib/dpkg/lock-frontend"

This typically means Ubuntu software updater is running in the background.  Stop the script, log out and back in again.  Or simply wait until the "Software Updater" UI appears on screen, and close it.

