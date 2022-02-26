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

func (e NumberBool) Bool() bool {
	return e == NumberBoolTrue
}

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

// ErrorNo - エラーNo
type ErrorNo string

const (
	ErrorUnspecified      ErrorNo = ""    // 未指定
	ErrorNoProblem        ErrorNo = "0"   // 問題なし
	ErrorNoData           ErrorNo = "1"   // データなし
	ErrorSessionInactive  ErrorNo = "2"   // 無効なセッション
	ErrorProgressedNumber ErrorNo = "6"   // 処理済みの送信通番
	ErrorExceedLimitTime  ErrorNo = "8"   // 送信日時からみたタイムアウト
	ErrorServiceOffline   ErrorNo = "9"   // サービス停止中
	ErrorBadRequest       ErrorNo = "-1"  // 引数エラー
	ErrorDatabaseAccess   ErrorNo = "-2"  // データベースへのアクセスエラー
	ErrorServerAccess     ErrorNo = "-3"  // サーバへのアクセスエラー
	ErrorSystemOffline    ErrorNo = "-12" // システム停止中
	ErrorOffHours         ErrorNo = "-62" // 情報提供時間外
)

// Exchange - 市場
type Exchange string

const (
	ExchangeUnspecified Exchange = ""   // 未指定
	ExchangeToushou     Exchange = "00" // 東証
)

// Side - 売買区分
type Side string

const (
	SideUnspecified Side = ""  // 未指定
	SideSell        Side = "1" // 売
	SideBuy         Side = "3" // 買
	SideDelivery    Side = "5" // 現渡
	SideReceipt     Side = "7" // 現引
)

// ExecutionTiming - 執行条件
type ExecutionTiming string

const (
	ExecutionTimingUnspecified ExecutionTiming = ""  // 未指定
	ExecutionTimingNormal      ExecutionTiming = "0" // 指定なし
	ExecutionTimingOpening     ExecutionTiming = "2" // 寄付
	ExecutionTimingClosing     ExecutionTiming = "4" // 引け
	ExecutionTimingFunari      ExecutionTiming = "6" // 不成
)

// TradeType - 現金信用区分
type TradeType string

const (
	TradeTypeUnspecified  TradeType = ""  // 未指定
	TradeTypeStock        TradeType = "0" // 現物
	TradeTypeSystemEntry  TradeType = "2" // 新規(制度信用6ヶ月)
	TradeTypeSystemExit   TradeType = "4" // 返済(制度信用6ヶ月)
	TradeTypeGeneralEntry TradeType = "6" // 新規(一般信用6ヶ月)
	TradeTypeGeneralExit  TradeType = "8" // 返済(一般信用6ヶ月)
)

// StopOrderType - 逆指値注文種別
type StopOrderType string

const (
	StopOrderTypeUnspecified   StopOrderType = ""  // 未指定
	StopOrderTypeNormal        StopOrderType = "0" // 通常
	StopOrderTypeStop          StopOrderType = "1" // 逆指値
	StopOrderTypeNormalAndStop StopOrderType = "2" // 通常 + 逆指値
)

// ExitOrderType - 建日種類
type ExitOrderType string

const (
	ExitOrderTypeUnspecified ExitOrderType = ""  // 未指定
	ExitOrderTypeUnused      ExitOrderType = " " // 未使用
	ExitOrderTypeSpecified   ExitOrderType = "1" // 個別指定
	ExitOrderTypeDayAsc      ExitOrderType = "2" // 建日順
	ExitOrderTypeProfitDesc  ExitOrderType = "3" // 単価益順
	ExitOrderTypeProfitAsc   ExitOrderType = "4" // 単価損順
)

// OrderInquiryStatus - 注文状態
type OrderInquiryStatus string

const (
	OrderInquiryStatusUnspecified OrderInquiryStatus = ""  // 未指定
	OrderInquiryStatusInOrder     OrderInquiryStatus = "1" // 未約定・注文中
	OrderInquiryStatusDone        OrderInquiryStatus = "2" // 全部約定
	OrderInquiryStatusPart        OrderInquiryStatus = "3" // 部分約定
	OrderInquiryStatusEditable    OrderInquiryStatus = "4" // 訂正取消可能な注文
	OrderInquiryStatusPartInOrder OrderInquiryStatus = "5" // 未約定 + 一部約定
)

// ExitTermType - 弁済区分
type ExitTermType string

const (
	ExitTermTypeUnspecified          ExitTermType = ""   // 未指定
	ExitTermTypeNoLimit              ExitTermType = "00" // 期限なし
	ExitTermTypeSystemMargin6m       ExitTermType = "26" // 制度信用6ヶ月
	ExitTermTypeSystemMarginNoLimit  ExitTermType = "29" // 制度信用無期限
	ExitTermTypeGeneralMargin6m      ExitTermType = "36" // 一般信用6ヶ月
	ExitTermTypeGeneralMarginNoLimit ExitTermType = "39" // 一般信用無期限
)

// ExecutionType - 注文値段区分
type ExecutionType string

const (
	ExecutionTypeUnspecified ExecutionType = ""  // 未指定
	ExecutionTypeUnused      ExecutionType = " " // 未使用
	ExecutionTypeMarket      ExecutionType = "1" // 成行
	ExecutionTypeLimit       ExecutionType = "2" // 指値
	ExecutionTypeHigher      ExecutionType = "3" // 親注文より高い
	ExecutionTypeLower       ExecutionType = "4" // 親注文より低い
)

// TriggerType - トリガータイプ
type TriggerType string

const (
	TriggerTypeUnspecified   TriggerType = ""  // 未指定
	TriggerTypeNoFired       TriggerType = "0" // 未発火
	TriggerTypeAuto          TriggerType = "1" // 自動
	TriggerTypeManualOrder   TriggerType = "2" // 手動発注
	TriggerTypeManualExpired TriggerType = "3" // 手動失効
)

// PartContractType - 内出来区分
type PartContractType string

const (
	PartContractTypeUnspecified PartContractType = ""  // 未指定
	PartContractTypeUnused      PartContractType = " " // 未使用
	PartContractTypePart        PartContractType = "2" // 分割約定
)

// OrderStatus - 状態コード
type OrderStatus string

const (
	OrderStatusUnspecified     OrderStatus = ""   // 未指定
	OrderStatusReceived        OrderStatus = "0"  // 受付未済
	OrderStatusInOrder         OrderStatus = "1"  // 未約定
	OrderStatusError           OrderStatus = "2"  // 受付エラー
	OrderStatusInCorrect       OrderStatus = "3"  // 訂正中
	OrderStatusCorrected       OrderStatus = "4"  // 訂正完了
	OrderStatusCorrectFailed   OrderStatus = "5"  // 訂正失敗
	OrderStatusInCancel        OrderStatus = "6"  // 取消中
	OrderStatusCanceled        OrderStatus = "7"  // 取消完了
	OrderStatusCancelFailed    OrderStatus = "8"  // 取消失敗
	OrderStatusPart            OrderStatus = "9"  // 一部約定
	OrderStatusDone            OrderStatus = "10" // 全部約定
	OrderStatusPartExpired     OrderStatus = "11" // 一部失効
	OrderStatusExpired         OrderStatus = "12" // 全部失効
	OrderStatusWait            OrderStatus = "13" // 発注待ち
	OrderStatusInvalid         OrderStatus = "14" // 無効
	OrderStatusTrigger         OrderStatus = "15" // 切替注文・逆指注文(切替中)
	OrderStatusTriggered       OrderStatus = "16" // 切替完了・逆指注文(未約定)
	OrderStatusTriggerFailed   OrderStatus = "17" // 切替失敗・逆指注文(失敗)
	OrderStatusCarryOverFailed OrderStatus = "19" // 繰越失効
	OrderStatusInOrderStop     OrderStatus = "50" // 逆指値発注中
)

// ContractStatus - 約定ステータス
type ContractStatus string

const (
	ContractStatusUnspecified ContractStatus = ""  // 未指定
	ContractStatusInOrder     ContractStatus = "0" // 未約定
	ContractStatusPart        ContractStatus = "1" // 部分約定
	ContractStatusDone        ContractStatus = "2" // 全部約定
	ContractStatusInContract  ContractStatus = "3" // 約定中
)

// CarryOverType - 繰越注文フラグ
type CarryOverType string

const (
	CarryOverTypeUnspecified CarryOverType = ""  // 未指定
	CarryOverTypeToday       CarryOverType = "0" // 当日
	CarryOverTypeCarry       CarryOverType = "1" // 繰越注文
	CarryOverTypeInvalid     CarryOverType = "2" // 無効
)

// CorrectCancelType - 訂正取消可否フラグ
type CorrectCancelType string

const (
	CorrectCancelTypeUnspecified CorrectCancelType = ""  // 未指定
	CorrectCancelTypeCorrectable CorrectCancelType = "0" // 訂正・取消可能
	CorrectCancelTypeCancelable  CorrectCancelType = "1" // 取消可能
	CorrectCancelTypeInvalid     CorrectCancelType = "2" // 訂正・取消不可
)

// Channel - チャネル
type Channel string

const (
	ChannelUnspecified Channel = ""  // 未指定
	ChannelMeet        Channel = "0" // 対面
	ChannelPC          Channel = "1" // PC
	ChannelCallCenter  Channel = "2" // コールセンター
	ChannelCallCenter2 Channel = "3" // コールセンター
	ChannelCallCenter3 Channel = "4" // コールセンター
	ChannelMobile      Channel = "5" // モバイル
	ChannelRich        Channel = "6" // リッチ
	ChannelSmartPhone  Channel = "7" // スマホ・タブレット
	ChannelIPadApp     Channel = "8" // iPadアプリ
	ChannelHost        Channel = "9" // HOST
)
