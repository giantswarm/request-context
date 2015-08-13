# request-context

The intention of this package is to have a context you can pass along methods
of your service objects. As functionality and complexity of systems increases
you need to provide more and more information from software layer to software
layer. Having a context object allows you to put arbitrary data into it and
pass one single context along your code flow.

### logging

Here we included a logger that can be provided with some context as well. Given
a context the logger marshals all information and print them additionally for
each line. In case you don't care, don't use this logger, or provide `nil` to
ignore additional payloads.
