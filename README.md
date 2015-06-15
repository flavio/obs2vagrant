# OBS2Vagrant

OBS2Vagrant is a simple web service that makes possible to turn a
[Open Build Service](http://openbuildservice.org/) into a simple
[Vagrant](https://www.vagrantup.com/) image catalog (like [HashiCorp's Atlas](https://atlas.hashicorp.com/)).

## The problem

It is possible to build Vagrant boxes using [KIWI](http://opensuse.github.io/kiwi/)
inside of the Open Build Service. The boxes are going to be automatically built
by OBS whenever one the packages used by them is changed or the KIWI spec file
is modified by the owner. The boxes are automatically versioned by OBS, meaning
their download url keeps changing over the time.

Referencing the download url of the box is not a good solution, because from
time to time it becames broken.

## The solution

As stated inside of [Vagrant's documentation](http://docs.vagrantup.com/v2/boxes/versioning.html):

  > Since Vagrant 1.5, boxes support versioning. This allows the people who make
  > boxes to push updates to the box, and the people who use the box have a
  > simple workflow for checking for updates, updating their boxes, and seeing
  > what has changed.

By default versioning works only when the box is registered on
[HashiCorp's Atlas](https://atlas.hashicorp.com/), however it's possible to
achieve the same functionality by pointing Vagrant to a specially crafted json
file. This is exatly what `obs2vagrant` does.

**Note well:** OBS hosts only latest version of the box.

## How obs2vagrant works

obs2vagrant is a small web application creating the special json files requested
by Vagrant to use box versioning. The json files are generated on the fly by
looking at the contents of a OBS project.

obs2vagrant does not use any database, it just uses some really simple conventions.

To obtain the special json file just make a `GET` request against:

  `/<server>/<project>/<repository>/<box_name>.json`

Where:
  * `server` is a unique ID identifying the OBS server instance hosting the box.
    This must be configured inside of obs2vagrant configuration file.
  * `project` is the name of the project on OBS.
  * `repository` is the name of the repository on OBS.
  * `box_name` is the name of the box, omit all version numbers.

The json file is created using the information obtained talking with
the OBS server using the [official OBS API](https://api.opensuse.org/apidocs/).

## Configuration file

obs2vagrant uses a simple configuration file to work, an example can be found
inside of this repository.

The configuration file is a json file like the following one:

```json
{
  "address" : "127.0.0.1",
  "port": 8080,
  "servers" : {
    "obs" : "http://download.opensuse.org/repositories/",
    "ibs" : "http://download.suse.de/ibs/"
  }
}
```

The following configuration defines two different server: "`obs`" and "`ibs`".
Each server must specify the root url for downloads.

## Example

Suppose you want to use the "`Base-SLES12-btrfs`" Vagrant box built on a private
OBS instance inside of the "`Devel:Docker:Images:KVM:SLE-12`" project. The project
builds the boxes inside of the repository named "`images`".

First of all you must add a server entry inside of your configuration file:
```json
"servers" : {
  "address" : "127.0.0.1",
  "port": 8080,
  "servers" : {
    "ibs" : "http://download.suse.de/ibs/"
  }
}
```

The special json file will be returned when visiting the following url:

`http://localhost:8080/ibs/Devel:Docker:Images:KVM:SLE-12/images/Base-SLES12-btrfs.json`


Your `Vagrantfile` will look something like:
```ruby
Vagrant.configure('2') do |config|

  config.vm.define :sles12 do |node|
    node.vm.box = 'Base-SLES12-btrfs'
    node.vm.box_url = 'http://localhost:8080/ibs/Devel:Docker:Images:KVM:SLE-12/images/Base-SLES12-btrfs.json'
    node.vm.box_check_update = true
  end
end
```

## Bulding & running

obs2vagrant is written using Go.

Just do:
`go install github.com/flavio/obs2vagrant`

To run it:

`obs2vagrant -c obs2vagrant.json`

## TODO

Code coverage.
