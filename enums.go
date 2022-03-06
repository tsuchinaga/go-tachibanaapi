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

// ResponseFormat - 応答データフォーマット指定
type ResponseFormat int

const (
	ResponseFormatUnspecified ResponseFormat = 0 // 未指定・標準
	ResponseFormatReadable    ResponseFormat = 1 // json文字列をタブや改行で成形された読みやすい形式
	ResponseFormatWrapped     ResponseFormat = 2 // {}の前後にのみ改行
	ResponseFormatWordKey     ResponseFormat = 4 // 引数項目番号ではなく引数項目での応答
)

// FeatureType - 機能種別
type FeatureType string

const (
	FeatureTypeUnspecified                 FeatureType = ""                                // 未指定
	FeatureTypeLoginRequest                FeatureType = "CLMAuthLoginRequest"             // ログインリクエスト
	FeatureTypeLoginResponse               FeatureType = "CLMAuthLoginAck"                 // ログインレスポンス
	FeatureTypeLogoutRequest               FeatureType = "CLMAuthLogoutRequest"            // ログアウトリクエスト
	FeatureTypeLogoutResponse              FeatureType = "CLMAuthLogoutAck"                // ログアウトレスポンス
	FeatureTypeNewOrder                    FeatureType = "CLMKabuNewOrder"                 // 新規注文
	FeatureTypeCorrectOrder                FeatureType = "CLMKabuCorrectOrder"             // 訂正注文
	FeatureTypeCancelOrder                 FeatureType = "CLMKabuCancelOrder"              // 取消注文
	FeatureTypeStockPositionList           FeatureType = "CLMGenbutuKabuList"              // 現物保有銘柄一覧
	FeatureTypeMarginPositionList          FeatureType = "CLMShinyouTategyokuList"         // 信用建玉一覧
	FeatureTypeStockWallet                 FeatureType = "CLMZanKaiKanougaku"              // 買余力
	FeatureTypeMarginWallet                FeatureType = "CLMZanShinkiKanoIjiritu"         // 建余力&本日維持率
	FeatureTypeSellable                    FeatureType = "CLMZanUriKanousuu"               // 売却可能数量 TODO これが何か調べる
	FeatureTypeOrderList                   FeatureType = "CLMOrderList"                    // 注文一覧
	FeatureTypeOrderListDetail             FeatureType = "CLMOrderListDetail"              // 注文約定一覧(詳細)
	FeatureTypeSummary                     FeatureType = "CLMZanKaiSummary"                // 可能額サマリー TODO これが何か調べる
	FeatureTypeSummaryRecord               FeatureType = "CLMZanKaiKanougakuSuii"          // 可能額推移 TODO これが何か調べる
	FeatureTypeStockEntryDetail            FeatureType = "CLMZanKaiGenbutuKaitukeSyousai"  // 現物株式買付可能額詳細 TODO これが何か調べる
	FeatureTypeMarginEntryDetail           FeatureType = "CLMZanKaiSinyouSinkidateSyousai" // 信用新規建て可能額詳細 TODO これが何か調べる
	FeatureTypeDepositRate                 FeatureType = "CLMZanRealHosyoukinRitu"         // リアル保証金率 TODO これが何か調べる
	FeatureTypeEventDownload               FeatureType = "CLMEventDownload"                // マスタ情報ダウンロード
	FeatureTypeEventSystemStatus           FeatureType = "CLMSystemStatus"                 // システムステータス
	FeatureTypeEventDate                   FeatureType = "CLMDateZyouhou"                  // 日付情報
	FeatureTypeEventTickGroup              FeatureType = "CLMYobine"                       // 呼値
	FeatureTypeEventOperationStatus        FeatureType = "CLMUnyouStatus"                  // 運用ステータス別状態
	FeatureTypeEventStockOperationStatus   FeatureType = "CLMUnyouStatusKabu"              // 運用ステータス(株式)
	FeatureTypeEventProductOperationStatus FeatureType = "CLMUnyouStatusHasei"             // 運用運用ステータス(派生)
	FeatureTypeStockMaster                 FeatureType = "CLMIssueMstKabu"                 // 株式銘柄マスタ
	FeatureTypeEventExchangeMaster         FeatureType = "CLMIssueSizyouMstKabu"           // 株式銘柄市場マスタ
	FeatureTypeEventStockRegulation        FeatureType = "CLMIssueSizyouKiseiKabu"         // 株式銘柄別・市場別規制
	FeatureTypeEventFutureMaster           FeatureType = "CLMIssueMstSak"                  // 先物銘柄マスタ
	FeatureTypeEventOptionMaster           FeatureType = "CLMIssueMstOp"                   // オプション銘柄マスタ TODO これが何か調べる
	FeatureTypeEventExchangeRegulation     FeatureType = "CLMIssueSizyouKiseiHasei"        // 派生銘柄別・市場別規制 TODO これが何か調べる
	FeatureTypeEventSubstitute             FeatureType = "CLMDaiyouKakeme"                 // 代用掛目 TODO これが何か調べる
	FeatureTypeEventDepositMaster          FeatureType = "CLMHosyoukinMst"                 // 保証金マスタ TODO これが何か調べる
	FeatureTypeEventErrorReason            FeatureType = "CLMOrderErrReason"               // 取引所エラー等理由コード TODO これが何か調べる
	FeatureTypeEventDownloadComplete       FeatureType = "CLMEventDownloadComplete"        // 初期ダウンロード終了通知
	FeatureTypeMasterData                  FeatureType = "CLMMfdsGetMasterData"            // マスタ情報
	FeatureTypeMarketPrice                 FeatureType = "CLMMfdsGetMarketPrice"           // 時価関連情報
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

// PrevCloseRatioType - 騰落率フラグ
type PrevCloseRatioType string

const (
	PrevCloseRatioTypeUnspecified PrevCloseRatioType = ""   // 未指定
	PrevCloseRatioTypeOver5       PrevCloseRatioType = "01" // +5.01%以上
	PrevCloseRatioTypeOver3       PrevCloseRatioType = "02" // +3.01%以上
	PrevCloseRatioTypeOver2       PrevCloseRatioType = "03" // +2.01%以上
	PrevCloseRatioTypeOver1       PrevCloseRatioType = "04" // +1.01%以上
	PrevCloseRatioTypeOver0       PrevCloseRatioType = "05" // +0.01%以上
	PrevCloseRatioTypeKeep        PrevCloseRatioType = "06" // 変化なし
	PrevCloseRatioTypeUnder0      PrevCloseRatioType = "07" // -0.01%以上
	PrevCloseRatioTypeUnder1      PrevCloseRatioType = "08" // -1.01%以上
	PrevCloseRatioTypeUnder2      PrevCloseRatioType = "09" // -2.01%以上
	PrevCloseRatioTypeUnder3      PrevCloseRatioType = "10" // -3.01%以上
	PrevCloseRatioTypeUnder5      PrevCloseRatioType = "11" // -5.01%以上
)

// ChangePriceType - 現値前値比較
type ChangePriceType string

const (
	ChangePriceTypeUnspecified       ChangePriceType = ""     // 未指定
	ChangePriceTypeNoChange          ChangePriceType = "0000" // 事象なし
	ChangePriceTypeEqual             ChangePriceType = "0056" // 現値=前値
	ChangePriceTypeRise              ChangePriceType = "0057" // 現値>前値
	ChangePriceTypeDown              ChangePriceType = "0058" // 現値<前値
	ChangePriceTypeOpenAfterStopping ChangePriceType = "0059" // 中断板寄後の始値
	ChangePriceTypeZaraba            ChangePriceType = "0060" // ザラバ引け
	ChangePriceTypeClose             ChangePriceType = "0061" // 板寄引け
	ChangePriceTypeCloseAtStopping   ChangePriceType = "0062" // 中断引け
	ChangePriceTypeStopping          ChangePriceType = "0068" // 売買停止引け
)

// IndicationPriceType - 気配値種類
type IndicationPriceType string

const (
	IndicationPriceTypeUnspecified              IndicationPriceType = ""     // 未指定
	IndicationPriceTypeNoChange                 IndicationPriceType = "0000" // 事象なし
	IndicationPriceTypeGeneral                  IndicationPriceType = "0101" // 一般気配
	IndicationPriceTypeSpecific                 IndicationPriceType = "0102" // 特別気配
	IndicationPriceTypeBeforeOpening            IndicationPriceType = "0107" // 寄前気配
	IndicationPriceTypeBeforeClosing            IndicationPriceType = "0108" // 停止前特別気配
	IndicationPriceTypeContinuance              IndicationPriceType = "0118" // 連続約定気配
	IndicationPriceTypeContinuanceBeforeClosing IndicationPriceType = "0119" // 停止前の連続約定気配
	IndicationPriceTypeMoving                   IndicationPriceType = "0120" // 一般気配、買上がり・売下がり中
)

// CurrentPriceType - 現在値種別
type CurrentPriceType string

const (
	CurrentPriceTypeUnspecified CurrentPriceType = ""     // 未指定
	CurrentPriceTypeStopHigh    CurrentPriceType = "0071" // ストップ高
	CurrentPriceTypeStopLow     CurrentPriceType = "0072" // ストップ安
)

// ApprovalLogin - ログイン許可区分
type ApprovalLogin string

const (
	ApprovalLoginUnspecified  ApprovalLogin = ""  // 未指定
	ApprovalLoginUnApproval   ApprovalLogin = "0" // 不許可
	ApprovalLoginApproval     ApprovalLogin = "1" // 許可
	ApprovalLoginOutOfService ApprovalLogin = "2" // 不許可(サービス時間外)
	ApprovalLoginTesting      ApprovalLogin = "9" // 管理者のみ(テスト中)
)

// SystemStatus - システム状態
type SystemStatus string

const (
	SystemStatusUnspecified SystemStatus = ""  // 未指定
	SystemStatusClosing     SystemStatus = "0" // 閉局
	SystemStatusOpening     SystemStatus = "1" // 開局
	SystemStatusPause       SystemStatus = "2" // 一時停止
)

// DayKey - 日付KEY
type DayKey string

const (
	DayKeyUnspecified DayKey = ""    // 未指定
	DayKeyToday       DayKey = "001" // 当日基準
	DayKeyNextDay     DayKey = "002" // 翌日基準
)

// TickGroupType - 呼値の単位番号
type TickGroupType string

const (
	TickGroupTypeUnspecified TickGroupType = ""    // 未指定
	TickGroupTypeStock1      TickGroupType = "101" // 株式1
	TickGroupTypeStock2      TickGroupType = "102" // 株式2
	TickGroupTypeStock3      TickGroupType = "103" // 株式3
	TickGroupTypeBond        TickGroupType = "201" // 債券1
	TickGroupTypeNK225       TickGroupType = "318" // 日経225先物
	TickGroupTypeNK225Mini   TickGroupType = "319" // 日経225mini先物
	TickGroupTypeNK225OP     TickGroupType = "418" // 日経225OP
)

// TaxFree - 非課税対象
type TaxFree string

const (
	TaxFreeUnspecified TaxFree = ""  // 未指定
	TaxFreeUnUsed      TaxFree = " " // 通常(無)
	TaxFreeValid       TaxFree = "1" // 非課税参加
)

// ExRightType - 権利落ちフラグ
type ExRightType string

const (
	ExRightTypeUnspecified              ExRightType = ""  // 未指定
	ExRightTypeNothing                  ExRightType = "0" // 権利落なし
	ExRightTypeStockSplit               ExRightType = "1" // 新株権利落ち
	ExRightTypeDividend                 ExRightType = "2" // 配当(中間)権利落ち
	ExRightTypeOther                    ExRightType = "4" // その他権利落ち
	ExRightTypeDividendAndOther         ExRightType = "5" // その他・配当(中間)権利落ち
	ExRightTypeStockSplitAndOther       ExRightType = "6" // 新株・その他権利落ち
	ExRightTypeStockSplitAndOtherMiddle ExRightType = "7" // 新株・その他(中間)権利落ち
)

// ListingType - 上場・入札C
type ListingType string

const (
	ListingTypeUnspecified  ListingType = ""  // 未指定
	ListingTypeUnUsed       ListingType = " " // 通常(無)
	ListingTypeNewest       ListingType = "1" // 上場1年未満銘柄
	ListingTypeGeneral      ListingType = "A" // 一般入札
	ListingTypeRight        ListingType = "B" // 権利入札
	ListingTypeOffer        ListingType = "C" // 公募入札
	ListingTypeSelling      ListingType = "D" // 売出し
	ListingTypeOpenBuy      ListingType = "E" // 公開買付
	ListingTypeTransmission ListingType = "F" // 媒介
)

// StopTradingType - 売買停止C
type StopTradingType string

const (
	StopTradingTypeUnspecified StopTradingType = ""  // 未指定
	StopTradingTypeUnUsed      StopTradingType = " " // 通常(無)
	StopTradingTypeRelease     StopTradingType = "0" // 解除
	StopTradingTypeStopping    StopTradingType = "9" // 停止中
)

// SettlementType - 決算C
type SettlementType string

const (
	SettlementTypeUnspecified     SettlementType = ""   // 未指定
	SettlementTypeCapitalIncrease SettlementType = "01" // 有償割当増資
	SettlementTypeSplit           SettlementType = "04" // 株式分割
	SettlementTypeAssignment      SettlementType = "05" // 無償割当
)

type MarginType string

const (
	MarginTypeUnspecified   MarginType = ""  // 未指定
	MarginTypeMarginTrading MarginType = "1" // 貸借銘柄
	MarginTypeSystem        MarginType = "2" // 信用制度銘柄
	MarginTypeGeneral       MarginType = "3" // 一般信用銘柄
)

// TradeRestriction - 取引禁止
type TradeRestriction string

const (
	TradeRestrictionUnspecified TradeRestriction = ""  // 未指定
	TradeRestrictionNormal      TradeRestriction = "0" // 通常(無)
	TradeRestrictionTrading     TradeRestriction = "1" // 取引禁止
	TradeRestrictionMarket      TradeRestriction = "2" // 成行禁止
	TradeRestrictionFraction    TradeRestriction = "3" // 端株禁止
)
