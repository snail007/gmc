# GMC CONTROLLER

## METHOD RULES

The rules about controller's method is below.

1. controller's method has a `__` suffix, method will be ignored by `router.Controller()`.

1. controller's method named `Before__()` is the construct method of controller, will be called before actual method call.

1. controller's method named `After__()` is the destruct method, will be called after actual method call.

## HELPER METHODS

GMC Controller defined some helper methods, help you to coding more easy.  

1. `Stop()`, call it, your code will exit current requested controller's method.

1. `Die()`, call it, your code will exit current requested controller's method, and prevent `After__` be called.

1. ``

## INSIDE METHODS

These inside controller's methods, don't call them in your code.

1. `MethodCallPre__()`, this method is be called before `Before__` to initialize base objects.

1. `MethodCallPost__()`, this method is be called after `After__` to do some ending works.

1. `Tr`, this is a i18n helper function, get details to read about `i18n/README.md`.

1. `SessionStart()` start a session, before you access session data, you must call `SessionStart` to start session.

1. `SessionDestroy()` destroy current session data.

## METHOD DEFINE

Your controller extends gmccontroller.Controller, so you must not be defined methods above in your controller.

## MEMBER DEFINE

