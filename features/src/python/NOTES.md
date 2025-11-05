## Notes

### .0 versions for old Python versions

Python <= 3.2 didn't use correct semver and therefore the .0 version was missing the .0 (eg. 3.1.0 was only released as 3.1).

Because of this, it is not possible to install an exact .0 version for those Python releases, it will take the newest patch release in that case.

Example:
* version 3.1 => would install 3.1.5
* version 3.1.0 => would unfortunately fail
