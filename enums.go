package tachibana

// Environment - 環境
type Environment string

const (
	EnvironmentUnspecified Environment = ""           // 未指定
	EnvironmentProduction  Environment = "production" // 本番環境
	EnvironmentDemo        Environment = "demo"       // デモ環境
)

// ApiVersion - APIのバージョン
type ApiVersion string

const (
	ApiVersionUnspecified ApiVersion = ""     // 未指定
	ApiVersionV4R1        ApiVersion = "v4r1" // V4R1
	ApiVersionV4R2        ApiVersion = "v4r2" // V4R2
	ApiVersionLatest                 = ApiVersionV4R2
)

// FeatureType - 機能種別
type FeatureType string

const (
	FeatureTypeUnspecified                FeatureType = ""                                // 未指定
	FeatureTypeLoginRequest               FeatureType = "CLMAuthLoginRequest"             // ログインリクエスト
	FeatureTypeLoginResponse              FeatureType = "CLMAuthLoginAck"                 // ログインレスポンス
	FeatureTypeLogoutRequest              FeatureType = "CLMAuthLogoutRequest"            // ログアウトリクエスト
	FeatureTypeLogoutResponse             FeatureType = "CLMAuthLogoutAck"                // ログアウトレスポンス
	FeatureTypeNewOrder                   FeatureType = "CLMKabuNewOrder"                 // 新規注文
	FeatureTypeCorrectOrder               FeatureType = "CLMKabuCorrectOrder"             // 訂正注文
	FeatureTypeCancelOrder                FeatureType = "CLMKabuCancelOrder"              // 取消注文
	FeatureTypeStockPositionList          FeatureType = "CLMGenbutuKabuList"              // 現物保有銘柄一覧
	FeatureTypeMarginPositionList         FeatureType = "CLMShinyouTategyokuList"         // 信用建玉一覧
	FeatureTypeStockWallet                FeatureType = "CLMZanKaiKanougaku"              // 買余力
	FeatureTypeMarginWallet               FeatureType = "CLMZanShinkiKanoIjiritu"         // 建余力&本日維持率
	FeatureTypeSellable                   FeatureType = "CLMZanUriKanousuu"               // 売却可能数量 TODO これが何か調べる
	FeatureTypeOrderList                  FeatureType = "CLMOrderList"                    // 注文一覧
	FeatureTypeOrderListDetail            FeatureType = "CLMOrderListDetail"              // 注文約定一覧(詳細)
	FeatureTypeSummary                    FeatureType = "CLMZanKaiSummary"                // 可能額サマリー TODO これが何か調べる
	FeatureTypeSummaryRecord              FeatureType = "CLMZanKaiKanougakuSuii"          // 可能額推移 TODO これが何か調べる
	FeatureTypeStockEntryDetail           FeatureType = "CLMZanKaiGenbutuKaitukeSyousai"  // 現物株式買付可能額詳細 TODO これが何か調べる
	FeatureTypeMarginEntryDetail          FeatureType = "CLMZanKaiSinyouSinkidateSyousai" // 信用新規建て可能額詳細 TODO これが何か調べる
	FeatureTypeDepositRate                FeatureType = "CLMZanRealHosyoukinRitu"         // リアル保証金率 TODO これが何か調べる
	FeatureTypeEventDownload              FeatureType = "CLMEventDownload"                // マスタ情報ダウンロード TODO これが何か調べる
	FeatureTypeEventSystemStatus          FeatureType = "CLMSystemStatus"                 // システムステータス TODO これが何か調べる
	FeatureTypeEventDate                  FeatureType = "CLMDateZyouhou"                  // 日付情報 TODO これが何か調べる
	FeatureTypeEventTick                  FeatureType = "CLMYobine"                       // 呼値 TODO これが何か調べる
	FeatureTypeEventOperationStatus       FeatureType = "CLMUnyouStatus"                  // 運用ステータス別状態 TODO これが何か調べる
	FeatureTypeEventStockOperationStatus  FeatureType = "CLMUnyouStatusKabu"              // 運用ステータス(株式) TODO これが何か調べる
	FeatureTypeEventMarginOperationStatus FeatureType = "CLMUnyouStatusHasei"             // 運用運用ステータス(派生) TODO これが何か調べる
	FeatureTypeEventStockMaster           FeatureType = "CLMIssueMstKabu"                 // 株式銘柄マスタ TODO これが何か調べる
	FeatureTypeEventExchangeMaster        FeatureType = "CLMIssueSizyouMstKabu"           // 株式銘柄市場マスタ TODO これが何か調べる
	FeatureTypeEventStockRegulation       FeatureType = "CLMIssueSizyouKiseiKabu"         // 株式銘柄別・市場別規制 TODO これが何か調べる
	FeatureTypeEventFutureMaster          FeatureType = "CLMIssueMstSak"                  // 先物銘柄マスタ TODO これが何か調べる
	FeatureTypeEventOptionMaster          FeatureType = "CLMIssueMstOp"                   // オプション銘柄マスタ TODO これが何か調べる
	FeatureTypeEventExchangeRegulation    FeatureType = "CLMIssueSizyouKiseiHasei"        // 派生銘柄別・市場別規制 TODO これが何か調べる
	FeatureTypeEventSubstitute            FeatureType = "CLMDaiyouKakeme"                 // 代用掛目 TODO これが何か調べる
	FeatureTypeEventDepositMaster         FeatureType = "CLMHosyoukinMst"                 // 保証金マスタ TODO これが何か調べる
	FeatureTypeEventErrorReason           FeatureType = "CLMOrderErrReason"               // 取引所エラー等理由コード TODO これが何か調べる
	FeatureTypeEventDownloadComplete      FeatureType = "CLMEventDownloadComplete"        // 初期ダウンロード終了通知 TODO これが何か調べる
)

// NumberBool - 数値表現されているbool
type NumberBool string

const (
	NumberBoolUnspecified NumberBool = ""  // 未指定
	NumberBoolFalse       NumberBool = "0" // false
	NumberBoolTrue        NumberBool = "1" // true
)

// AccountType - 口座種別
type AccountType string

const (
	AccountTypeUnspecified AccountType = ""  // 未指定
	AccountTypeSpecific    AccountType = "1" // 特定
	AccountTypeGeneral     AccountType = "3" // 一般
	AccountTypeNISA        AccountType = "5" // NISA
)

// SpecificAccountType - 特定口座区分
type SpecificAccountType string

const (
	SpecificAccountTypeUnspecified SpecificAccountType = ""  // 未指定
	SpecificAccountTypeGeneral     SpecificAccountType = "0" // 一般
	SpecificAccountTypeNothing     SpecificAccountType = "1" // 特定源泉徴収なし
	SpecificAccountTypeWithholding SpecificAccountType = "2" // 特定源泉徴収あり
)

// ErrType - エラーNo
type ErrType string

const (
	ErrTypeUnspecified      ErrType = ""    // 未指定
	ErrTypeNoProblem        ErrType = "0"   // 問題なし
	ErrTypeNoData           ErrType = "1"   // データなし
	ErrTypeSessionInactive  ErrType = "2"   // 無効なセッション
	ErrTypeProgressedNumber ErrType = "6"   // 処理済みの送信通番
	ErrTypeExceedLimitTime  ErrType = "8"   // 送信日時からみたタイムアウト
	ErrTypeServiceOffline   ErrType = "9"   // サービス停止中
	ErrTypeBadRequest       ErrType = "-1"  // 引数エラー
	ErrTypeDatabaseAccess   ErrType = "-2"  // データベースへのアクセスエラー
	ErrTypeServerAccess     ErrType = "-3"  // サーバへのアクセスエラー
	ErrTypeSystemOffline    ErrType = "-12" // システム停止中
	ErrTypeOffHours         ErrType = "-62" // 情報提供時間外
)
