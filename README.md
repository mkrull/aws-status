This is a scraping check of AWS service availability.

To build and run the project a docker client has to be installed on the
machines building and running the images.

Install:

    $ git clone https://github.com/mkrull/aws-status.git

Build:

    $ cd aws-status
    $ make

Building results in a docker image that can be run with `make` as well:

    $ make run

The output is a list of AWS services of the four regions with their current
operational status.

In case a service is degraded a notifier function is run.

Currrently this means the status is logged at the end of the service listing.
