# Sprig Function Documentation

The Sprig library provides over 148 template functions for Go's template language.

- [String Functions](strings.md): `trim`, `randAlpha`, `plural`, etc.
  - [String List Functions](string_slice.md): `splitList`, `sortAlpha`, etc.
- [Math Functions](math.md): `add`, `max`, `mul`, etc.
  - [Integer Slice Functions](integer_slice.md): `until`, `untilStep`
- [Date Functions](date.md): `now`, `date`, etc.
- [Defaults Functions](defaults.md): `default`, `empty`, `coalesce`, `toJSON`, `toPrettyJSON`, `toRawJSON`, `ternary`
- [Encoding Functions](encoding.md): `b64enc`, `b64dec`, etc.
- [Lists and List Functions](lists.md): `first`, `uniq`, etc.
- [Type Conversion Functions](conversion.md): `atoi`, `int64`, `toString`, etc.
- [File Path Functions](paths.md): `base`, `dir`, `ext`, `clean`, `isAbs`
- [Flow Control Functions](flow_control.md): `fail`
- Advanced Functions
  - [UUID Functions](uuid.md): `uuidv4`
  - [OS Functions](os.md): `env`, `expandenv`
  - [Reflection](reflection.md): `typeOf`, `kindIs`, `typeIsLike`, etc.
  - [Cryptographic and Security Functions](crypto.md): `derivePassword`, `sha256sum`, `genPrivateKey`, etc.

