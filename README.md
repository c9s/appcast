Appcast XML package for Go
==========================

```
appcast -append  \
	-title="1 bug fix" \
	-description="release note" \
	-file tests/tests.zip \
	-url http://yourhost.com \
	-version 109 \
	-dsaSignature="blah" \
	-versionShortString="1.4" \
		tests/appcast.xml
```

To parse namespace attributes for `sparkle:`, you need a patch for encoding/xml,
please download the patch from the below URL:

	https://codereview.appspot.com/download/issue7350048_10002.diff

License
-------
Public Domain



