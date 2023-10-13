Your first test
===============

.. warning::

  New Documentation https://yandex.github.io/pandora/
;
You can use Pandora alone or use it with `Yandex.Tank`_ as a test runner and
`Overload`_ as a result viewer. In the second case Pandora's configuration is the same, but you will embed it into Yandex.Tank's config.

Config file
-----------

Pandora supports config files in `YAML`_ format. Create a new file named ``load.yaml`` and add following lines in your favourite editor:

.. code-block:: yaml

  pools:
    - id: HTTP pool                    # pool name (for your choice)
      gun:
        type: http                     # gun type
        target: example.com:80         # gun target
      ammo:
        type: uri                      # ammo format
        file: ./ammo.uri               # ammo File
      result:
        type: phout                    # report format (phout is compatible with Yandex.Tank)
        destination: ./phout.log       # report file name

      rps:                             # shooting schedule
        type: line                     # linear growth
        from: 1                        # from 1 response per second
        to: 5                          # to 5 responses per second
        duration: 60s                  # for 60 seconds

      startup:                         # instances startup schedule
        type: once                     # start 10 instances
        times: 10


Visit `Configuration`_ page for more details.


Run your tests:


.. code-block:: bash

  pandora load.yaml


The results are in ``phout.log``. Use `Yandex.Tank`_
and `Overload`_ to plot them.

References
----------

.. target-notes::

.. _`Overload`: https://overload.yandex.net
.. _`Yandex.Tank`: https://yandextank.readthedocs.io/en/latest/core_and_modules.html#pandora
.. _`YAML`: https://en.wikipedia.org/wiki/YAML
.. _`Configuration`: configuration.html
