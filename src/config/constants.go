package config

// Network configuration constants
const (
	DefaultScheme = "ws"
	DefaultPath   = "/chat"
	DefaultPort   = "1337"
	SecureScheme  = "wss"
)

// Encryption constants
const (
	RSAKeySize   = 2048
	AESKeySize   = 32 // 256 bits
	AESBlockSize = 16 // 128 bits
)

// UI constants
const (
	MaxMessages         = 100
	DefaultInputWidth   = 50
	DefaultViewportW    = 50
	DefaultViewportH    = 20
	DefaultOnlineWidth  = 20
	DefaultOnlineHeight = 20
)

// Color constants
const (
	DefaultColorRange = 256
	HelpColorCommand  = "help"
)

// Command constants
const (
	JoinCommand  = "/join"
	LeaveCommand = "/leave"
	ExitCommand  = "/exit"
	ClearCommand = "/clear"
)

// Error messages
const (
	ErrEmptyConfig      = "configuration is empty"
	ErrEmptyHost        = "please provide IP and port of server and try again, use -h for help"
	ErrColorList        = "color list requested"
	ErrConnectionFailed = "failed establishing connection"
	ErrHandshakeFailed  = "failed sending initial handshake"
)
