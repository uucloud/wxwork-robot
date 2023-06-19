package sender

type SendResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type SendReq struct {
	CommonSend

	Text         *TextType     `json:"text,omitempty"`
	Markdown     *MarkdownType `json:"markdown,omitempty"`
	Image        *ImageType    `json:"image,omitempty"`
	News         *NewsType     `json:"news,omitempty"`
	File         *FileType     `json:"file,omitempty"`
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
}

type CommonSend struct {
	//会话id，支持最多传100个，用‘|’分隔。目前仅支持群聊会话，通过消息回调获得
	ChatID string `json:"chatid,omitempty"`

	//text/markdown/image/news/file/template_card
	MsgType string `json:"msgtype"`

	DelayUpdate uint32 `json:"delay_update,omitempty"`
}

type TextType struct {
	Content string `json:"content"`

	//userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，开发者可以通过回调事件中获取userid。
	//如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedList []string `json:"mentioned_list,omitempty"`

	//手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type MarkdownType struct {
	//markdown内容，最长不超过4096个字节，必须是utf8编码。点击查看目前支持的markdown语法
	//特殊的，content中支持<@userid>的@人语法，开发者可以通过回调事件中获取userid
	Content string `json:"content"`
}

type ImageType struct {
	//图片内容的base64编码
	Base64 string `json:"base64"`
	//图片内容（base64编码前）的md5值
	Md5 string `json:"md5"`
}

type NewsType struct {
	//图文消息，一个图文消息支持1到8条图文
	Articles []Article `json:"articles"`
}

type Article struct {
	//标题，不超过128个字节，超过会自动截断
	Title string `json:"title"`

	//描述，不超过512个字节，超过会自动截断
	Description string `json:"description,omitempty"`

	//点击后跳转的链接
	URL string `json:"url"`

	//图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
	Picurl string `json:"picurl,omitempty"`
}

type FileType struct {
	//文件id，通过机器人的文件上传接口获取
	MediaID string `json:"media_id"`
}

type TemplateCard struct {
	//模版卡片的模版类型，文本通知模版卡片的类型为text_notice
	CardType string `json:"card_type"`

	//卡片来源样式信息，不需要来源样式可不填写
	Source *TemplateCardSource `json:"source,omitempty"`

	//卡片右上角更多操作按钮
	ActionMenu *TemplateCardActionMenu `json:"action_menu,omitempty"`

	//模版卡片的主要内容，包括一级标题和标题辅助信息
	MainTitle *TemplateCardMainTitle `json:"main_title,omitempty"`

	//关键数据样式，建议不与引用样式共用
	EmphasisContent *TemplateCardEmphasisContent `json:"emphasis_content,omitempty"`

	//引用文献样式，建议不与关键数据共用
	QuoteArea *TemplateCardQuoteArea `json:"quote_area,omitempty"`

	//二级普通文本，建议不超过112个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写
	SubTitleText string `json:"sub_title_text,omitempty"`

	//二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6
	HorizontalContentList []*TemplateCardHorizontal `json:"horizontal_content_list,omitempty"`

	//跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3
	JumpList []*TemplateJumpList `json:"jump_list,omitempty"`

	//整体卡片的点击跳转事件，text_notice模版卡片中该字段为必填项
	CardAction *TemplateCardAction `json:"card_action"`

	//图片样式，news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填
	CardImage *TemplateCardImage `json:"card_image,omitempty"`

	//左图右文样式
	ImageText *TemplateCardImageText `json:"image_text,omitempty"`

	//卡片二级垂直内容，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过4
	VerticalContent []*TemplateCardVerticalContent `json:"vertical_content,omitempty"`

	//下拉式的选择器
	ButtonSelection *TemplateCardButtonSelection `json:"button_selection,omitempty"`

	//按钮列表，列表长度不超过6
	ButtonList []*TemplateCardButtonList `json:"button_list,omitempty"`

	//下拉式的选择器列表，multiple_interaction类型的卡片该字段不可为空，一个消息最多支持 3 个选择器
	SelectList []*TemplateCardSelectList `json:"select_list,omitempty"`

	//提交按钮样式
	SubmitButton *TemplateCardSubmitButton `json:"submit_button,omitempty"`

	//任务id，当文本通知模版卡片有action_menu字段的时候，该字段必填。
	//同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。
	//任务id将会在相应的回调事件中返回
	TaskID string `json:"task_id,omitempty"`
}

type TemplateCardSource struct {
	//来源图片的url
	IconURL string `json:"icon_url,omitempty"`

	//来源图片的描述，建议不超过13个字
	Desc string `json:"desc,omitempty"`

	//来源文字的颜色，目前支持：0(默认) 灰色，1 黑色，2 红色，3 绿色
	DescColor int `json:"desc_color,omitempty"`
}

type TemplateCardActionMenu struct {
	//更多操作界面的描述
	Desc string `json:"desc"`

	//操作列表，列表长度取值范围为 [1, 3]
	ActionList []struct {
		//操作的描述文案
		Text string `json:"text"`
		//操作key值，用户点击后，会产生回调事件将本参数作为EventKey返回，回调事件会带上该key值，最长支持1024字节，不可重复
		Key string `json:"key"`
	} `json:"action_list"`
}

type TemplateCardMainTitle struct {
	//一级标题，建议不超过26个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写
	Title string `json:"title,omitempty"`
	//标题辅助信息，建议不超过30个字
	Desc string `json:"desc,omitempty"`
}

type TemplateCardEmphasisContent struct {
	//关键数据样式的数据内容，建议不超过10个字
	Title string `json:"title,omitempty"`
	//关键数据样式的数据描述内容，建议不超过15个字
	Desc string `json:"desc,omitempty"`
}

type TemplateCardQuoteArea struct {
	//引用文献样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type,omitempty"`
	//点击跳转的url，quote_area.type是1时必填
	URL string `json:"url,omitempty"`
	//点击跳转的小程序的appid，必须是与当前应用关联的小程序，quote_area.type是2时必填
	AppID string `json:"appid,omitempty"`
	//点击跳转的小程序的pagepath，quote_area.type是2时选填
	Pagepath string `json:"pagepath,omitempty"`
	//引用文献样式的标题
	Title string `json:"title,omitempty"`
	//引用文献样式的引用文案
	QuoteText string `json:"quote_text,omitempty"`
}

type TemplateCardHorizontal struct {
	//链接类型，0或不填代表是普通文本，1 代表跳转url，2 代表下载附件，3 代表点击跳转成员详情
	Type int `json:"type,omitempty"`
	//二级标题，建议不超过5个字
	Keyname string `json:"keyname"`
	//二级文本，如果horizontal_content_list.type是2，该字段代表文件名称（要包含文件类型），建议不超过26个字
	Value string `json:"value,omitempty"`
	//链接跳转的url，horizontal_content_list.type是1时必填
	URL string `json:"url,omitempty"`
	//附件的media_id，horizontal_content_list.type是2时必填
	MediaID string `json:"media_id,omitempty"`
	//成员详情的userid，horizontal_content_list.type是3时必填
	UserID string `json:"user_id,omitempty"`
}

type TemplateJumpList struct {
	//跳转链接类型，0或不填代表不是链接，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type,omitempty"`
	//跳转链接样式的文案内容，建议不超过13个字
	Title string `json:"title"`
	//跳转链接的url，jump_list.type是1时必填
	URL string `json:"url,omitempty"`
	//跳转链接的小程序的appid，jump_list.type是2时必填
	AppID string `json:"appid,omitempty"`
	//跳转链接的小程序的pagepath，jump_list.type是2时选填
	Pagepath string `json:"pagepath,omitempty"`
}

type TemplateCardAction struct {
	//卡片跳转类型，0或不填代表不是链接，1 代表跳转url，2 代表打开小程序。text_notice模版卡片中该字段取值范围为[1,2]
	Type int `json:"type"`
	//跳转事件的url，card_action.type是1时必填
	URL string `json:"url,omitempty"`
	//跳转事件的小程序的appid，card_action.type是2时必填
	AppID string `json:"appid,omitempty"`
	//跳转事件的小程序的pagepath，card_action.type是2时选填
	Pagepath string `json:"pagepath,omitempty"`
}

type TemplateCardImage struct {
	//图片的url
	URL string `json:"url"`
	//图片的宽高比，宽高比要小于2.25，大于1.3，不填该参数默认1.3
	AspectRatio float64 `json:"aspect_ratio,omitempty"`
}

type TemplateCardImageText struct {
	//左图右文样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type,omitempty"`
	//点击跳转的url，image_text_area.type是1时必填
	URL string `json:"url,omitempty"`
	//点击跳转的小程序的appid，必须是与当前应用关联的小程序，image_text_area.type是2时必填
	AppID string `json:"appid,omitempty"`
	//点击跳转的小程序的pagepath，image_text_area.type是2时选填
	Pagepath string `json:"pagepath,omitempty"`
	//左图右文样式的标题
	Title string `json:"title,omitempty"`
	//左图右文样式的描述
	Desc string `json:"desc,omitempty"`
	//左图右文样式的图片url
	ImageURL string `json:"image_url"`
}

type TemplateCardVerticalContent struct {
	//卡片二级标题，建议不超过26个字
	Title string `json:"title"`
	//二级普通文本，建议不超过112个字
	Desc string `json:"desc,omitempty"`
}

type TemplateCardButtonSelection struct {
	//下拉式的选择器的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节
	QuestionKey string `json:"question_key"`
	//下拉式的选择器左边的标题
	Title string `json:"title,omitempty"`
	//下拉式的选择器是否不可选，false为可选，true为不可选。仅在更新模版卡片的时候该字段有效
	Disable bool `json:"disable,omitempty"`
	//选项列表，下拉选项不超过 10 个，最少1个
	OptionList struct {
		//下拉式的选择器选项的id，用户提交后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复
		ID string `json:"id"`
		//下拉式的选择器选项的文案，建议不超过16个字
		Text string `json:"text"`
	} `json:"option_list"`
	//默认选定的id，不填或错填默认第一个
	SelectedID string `json:"selected_id"`
}

type TemplateCardButtonList struct {
	//按钮文案，建议不超过10个字
	Text string `json:"text"`
	//按钮样式，目前可填1~4，不填或错填默认1, 1:蓝底白字 2:蓝字白底 3:红字白底 4:黑字白底
	Style int `json:"style,omitempty"`
	//按钮key值，用户点击后，会产生回调事件将本参数作为event_key返回，最长支持1024字节，不可重复
	Key string `json:"key"`
}

type TemplateCardSelectList struct {
	//下拉式的选择器题目的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节，不可重复
	QuestionKey string `json:"question_key"`

	//选择器的标题，建议不超过13个字
	Title string `json:"title,omitempty"`

	//下拉式的选择器是否不可选，false为可选，true为不可选。仅在更新模版卡片的时候该字段有效
	Disable bool `json:"disable,omitempty"`

	//默认选定的id，不填或错填默认第一个
	SelectedID string `json:"selected_id,omitempty"`

	//选项列表，下拉选项不超过 10 个，最少1个
	OptionList []TemplateCardSelectOption `json:"option_list"`
}

type TemplateCardSelectOption struct {
	//下拉式的选择器选项的id，用户提交选项后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复
	ID string `json:"id"`
	//下拉式的选择器选项的文案，建议不超过10个字
	Text string `json:"text"`
}

type TemplateCardSubmitButton struct {
	//按钮文案，建议不超过10个字
	Text string `json:"text"`
	//提交按钮的key，会产生回调事件将本参数作为EventKey返回，最长支持1024字节
	Key string `json:"key"`
}
