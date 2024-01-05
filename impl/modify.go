package impl

const (
	PUBLIC        = 0x00000001 /* visible to everyone */
	PRIVATE       = 0x00000002 /* visible only to the defining class */
	PROTECTED     = 0x00000004 /* visible to subclasses */
	STATIC        = 0x00000008 /* instance variable is static */
	FINAL         = 0x00000010 /* no further subclassing, overriding */
	SYNCHRONIZED  = 0x00000020 /* wrap method call in monitor lock */
	VOLATILE      = 0x00000040 /* can cache in registers */
	BRIDGE        = 0x00000040 /* Bridge method generated by compiler */
	TRANSIENT     = 0x00000080 /* not persistant */
	VARARGS       = 0x00000080 /* Method accepts var. args*/
	NATIVE        = 0x00000100 /* implemented in C */
	INTERFACE     = 0x00000200 /* class is an interface */
	ABSTRACT      = 0x00000400 /* no definition provided */
	ENUM_CONSTANT = 0x00004000 /* enum constant field*/
	SYNTHETIC     = 0xf0000000 /* not in source code */
)