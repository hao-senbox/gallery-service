package constants

var (
	GalleryLanguages = map[string]Language{
		"en": EnglishLanguageConfig,
		"vi": VietnameseLanguageConfig,
		"de": GermanLanguageConfig,
		"fr": FrenchLanguageConfig,
		"es": SpanishLanguageConfig,
	}
	Components = map[string]Component{
		"time-circle":      TimerCircleComponent,
		"time-linear":      TimerLinearComponent,
		"time-linear-full": TimerLinearFullComponent,
		"time-power":       TimerPowerComponent,
		"todo":             TodoComponent,
		"schedule":         ScheduleComponent,
		"choice-board":     ChoiceBoardComponent,
		"activity":         ActivityComponent,
		"video-player":     VideoPlayerComponent,
		"album":            AlbumComponent,
		"economy":          EconomyComponent,
		"chat":             ChatComponent,
		"alarm":            AlarmComponent,
		"sentence":         SentenceComponent,
	}
)

type Language string

const (
	EnglishLanguageConfig    Language = "English"
	VietnameseLanguageConfig Language = "Vietnamese"
	GermanLanguageConfig     Language = "German"
	FrenchLanguageConfig     Language = "French"
	SpanishLanguageConfig    Language = "Spanish"
)

func (l Language) String() string {
	return string(l)
}

type Component string

const (
	TimerCircleComponent     Component = "Timer Circle"
	TimerLinearComponent     Component = "Timer Linear"
	TimerLinearFullComponent Component = "Timer Linear Full"
	TimerPowerComponent      Component = "Timer Power"
	TodoComponent            Component = "Todo"
	ScheduleComponent        Component = "Schedule"
	ChoiceBoardComponent     Component = "Choice Board"
	ActivityComponent        Component = "Activity (SBT)"
	VideoPlayerComponent     Component = "Video Player"
	AlbumComponent           Component = "Album"
	EconomyComponent         Component = "Economy"
	ChatComponent            Component = "Chat"
	AlarmComponent           Component = "Alarm"
	SentenceComponent        Component = "Sentence"
)

func (l Component) String() string {
	return string(l)
}
