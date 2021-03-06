## RHEL Satellite (5) and issues with it

1. **Package leaking on frozen channels**
    When we are running with frozen channels and have introduced new packages
    and/or security patches we experience that the redhat satellite introduces
    new patch releases and add packages which we explicitly held back. This
    even though they are not a required package by the patches we actually
    want.

2. **Unlocked channels**
    This is what we are running with now. And is to my knowledge the prefered
    way to run. This results in newer packages dropping in every now and then
    and a slight derivation between installed servers. Pretty much the same as
    having leaking packages but without the false security and that it's not
    just a few that are leaking.

3. **Limitations from upstream**
    The Satellite release we are running is a EIS standard and is located in
    Kista, all sites are sharing this and only have a simple proxy locally. So
    if we wish to have the new satellite it has to be a global way of working
    change. (Unless we can setup one in parallel).

4. **Ways forward**
    The new satellite is supposed to be rewritten from scratch. Introducing
    propper revision handling of package channels and pre/post-scripts. As well
    as integration with Puppet

## Repository
### Full package sync
```
/upstream
    /rhel
        /7
            /x86_64
                /base
                /debuginfo
                /extras
                /optional
                /supplementary
```

These folders only contain 1 release of each package.

These products only have the packages required for their specific purpose

```
/products
    /common-rh7
        /packages            # packagelist with all packages in each commit
            /1.0.1
            /1.1.0
            /latest -> ./1.1.0
            /stable -> ./1.0.1
        /1.0.1                                  # (generated)
            /repodata                           # (generated by createrepo)
        /1.1.0                                  # (generated)
            /repodata                           # (generated by createrepo)
        /latest -> ./1.1.0                      # (generated)
        /stable -> ./1.0.1                      # (generated)
    /minimal-rh7
        /packages            # packagelist with all packages in each commit
            /1.1.0
            /1.2.0
            /latest -> ./1.2.0
            /stable -> ./1.1.0
        /1.0.1                                  # (generated)
            /repodata                           # (generated by createrepo)
        /1.1.0                                  # (generated)
            /repodata                           # (generated by createrepo)
        /1.2.0                                  # (generated)
            /repodata                           # (generated by createrepo)
        /latest -> ./1.2.0                      # (generated)
        /stable -> ./1.1.0                      # (generated)
```
All except upstream will simply be symlinks to the packages.

### Tag explanation
The tag will adhere to a bastard version of [Semantic Versioning](http://semver.org/)
where the 0.Y.Z releases will be while trimming the package install,
while in production: X.Y.Z
X = Removal of package(s) (ie breaking api)
Y = Introduce new package(s) (adding a feature)
Z = Patch package(s) (bugfix)

Introducing major version uppgrades of packages might merit an increment on X or Y, use your own judgement for this, though this will be rare in RHEL at least.


## Kickstart
When Kickstart is run it requires Repository to be completed and have generated the required url's which point to the package list.
### minimal-rh7.yaml
```
os: rhel
os_release: 7
product: minimal-rh7
product_release: stable
```

```
/proj/selnhubadm/Satellite
    /kickstart
        /common-rh7-stable                       # (generated)
        /common-rh7-latest                       # (generated)
        /minimal-rh7-stable                      # (generated)
        /minimal-rh7-latest                      # (generated)
```

## Prerequsite
rpm sync from upstream, this will not be handled in this module in the forseeable future.

dispatch stuff to pxe server, see example folder for how this can be done.
