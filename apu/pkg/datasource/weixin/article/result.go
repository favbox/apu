package article

// 微信读书的文章列表
type (
	BookArticle struct {
		Time    int    `json:"time"`
		Avatar  string `json:"avatar"`
		MpName  string `json:"mp_name"`
		Title   string `json:"title"`
		Content string `json:"content"`
		DocUrl  string `json:"doc_url"`
		PicUrl  string `json:"pic_url"`
	}

	ReviewInfo struct {
		Review struct {
			BelongBookId string       `json:"belongBookId"`
			MpInfo       *BookArticle `json:"mpInfo"`
		} `json:"review"`
		ReviewId string `json:"reviewId"`
	}

	ArticlesResult struct {
		Errcode int    `json:"errcode"`
		Errlog  string `json:"errlog"`
		Errmsg  string `json:"errmsg"`

		ClearAll int           `json:"clearAll"`
		Reviews  []*ReviewInfo `json:"reviews"`
		SyncKey  int           `json:"synckey"`
	}
)

// 文章统计
type (
	Stat struct {
		CollectNum      int  `json:"collect_num"`
		FriendLikeNum   int  `json:"friend_like_num"`
		IsLogin         bool `json:"is_login"`
		LikeDisabled    bool `json:"like_disabled"`
		LikeNum         int  `json:"like_num"`
		Liked           bool `json:"liked"`
		OldLikeNum      int  `json:"old_like_num"`
		OldLiked        bool `json:"old_liked"`
		OldLikedBefore  int  `json:"old_liked_before"`
		Prompted        int  `json:"prompted"`
		ReadNum         int  `json:"read_num"`
		RealReadNum     int  `json:"real_read_num"`
		Ret             int  `json:"ret"`
		ShareNum        int  `json:"share_num"`
		Show            bool `json:"show"`
		ShowGray        int  `json:"show_gray"`
		ShowLike        int  `json:"show_like"`
		ShowLikeGray    int  `json:"show_like_gray"`
		ShowOldLike     int  `json:"show_old_like"`
		ShowOldLikeGray int  `json:"show_old_like_gray"`
		ShowRead        int  `json:"show_read"`
		Style           int  `json:"style"`
		Version         int  `json:"version"`
		VideoPv         int  `json:"video_pv"`
		VideoUv         int  `json:"video_uv"`
	}

	StatResult struct {
		AdvertisementInfo []any  `json:"advertisement_info"`
		Appid             string `json:"appid"`
		AppmsgAlbumVideos []any  `json:"appmsg_album_videos"`
		Appmsgact         struct {
			FavoriteBefore int `json:"favorite_before"`
			FollowBefore   int `json:"follow_before"`
			OldLikedBefore int `json:"old_liked_before"`
			PayBefore      int `json:"pay_before"`
			RewardBefore   int `json:"reward_before"`
			SeenBefore     int `json:"seen_before"`
			ShareBefore    int `json:"share_before"`
		} `json:"appmsgact"`
		ArticleStat *Stat `json:"appmsgstat"`
		BaseResp    struct {
			ExportkeyToken string `json:"exportkey_token"`
			Ret            int    `json:"ret"`
		} `json:"base_resp"`
		BizfileRet                      int `json:"bizfile_ret"`
		CloseRelatedArticle             int `json:"close_related_article"`
		DistanceToGetRelatedArticleData int `json:"distance_to_get_related_article_data"`
		FavoriteFlag                    struct {
			Show     int `json:"show"`
			ShowGray int `json:"show_gray"`
		} `json:"favorite_flag"`
		FriendSubscribeCount int   `json:"friend_subscribe_count"`
		HitBizrecommend      int   `json:"hit_bizrecommend"`
		IsFans               int   `json:"is_fans"`
		LinkComponentList    []any `json:"link_component_list"`
		MoreReadList         []any `json:"more_read_list"`
		OriginalArticleCount int   `json:"original_article_count"`
		PublicTagInfo        struct {
			Tags []any `json:"tags"`
		} `json:"public_tag_info"`
		RelatedArticleFastClose int   `json:"related_article_fast_close"`
		RelatedArticleUnderAd   int   `json:"related_article_under_ad"`
		RelatedTagVideo         []any `json:"related_tag_video"`
		RewardHeadImgInfos      []any `json:"reward_head_img_infos"`
		RewardHeadImgs          []any `json:"reward_head_imgs"`
		SecControl              struct {
			AdViolationMiddlePage int `json:"ad_violation_middle_page"`
		} `json:"sec_control"`
		ShareFlag struct {
			Show     int `json:"show"`
			ShowGray int `json:"show_gray"`
		} `json:"share_flag"`
		ShowBizBanner         int    `json:"show_biz_banner"`
		ShowRelatedArticle    int    `json:"show_related_article"`
		ShowRelatedSearchWord int    `json:"show_related_search_word"`
		TestFlag              int    `json:"test_flag"`
		VideoContinueFlag     int    `json:"video_continue_flag"`
		WapExportToken        string `json:"wap_export_token"`
	}
)
