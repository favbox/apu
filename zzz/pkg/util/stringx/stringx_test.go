package stringx

import (
	"fmt"
	"regexp"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	n, err := Parse[int32]("123")
	assert.Nil(t, err)
	assert.Equal(t, int32(123), n)

	b, err := Parse[bool]("1")
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	n2 := MustNumber[int]("123")
	assert.Equal(t, 123, n2)

	n3 := MustNumber[int]("138.688px")
	assert.Equal(t, 138, n3)
}

// 提取 style 属性中的 width 值
func ExtractWidthFromStyle(style string) (string, bool) {
	// 正则表达式匹配 width 属性，考虑可能存在的空格
	re := regexp.MustCompile(`(?i)width\s*:\s*(\d+)(px)?`)
	match := re.FindStringSubmatch(style)

	if len(match) < 2 {
		return "", false
	}

	// 返回第一个捕获组，即数字部分
	return match[1], true
}

func TestExtractWidth(t *testing.T) {
	styles := []string{
		"width: 200px; height: 100px;",
		"width :150px; height: 100px;",
		"width:150; height: 80;",
		"display:block; width: 300px;",
		"WIDTH: 150px;", // 测试大小写不敏感
	}

	for _, style := range styles {
		width, exists := ExtractWidthFromStyle(style)
		if exists {
			fmt.Printf("Width extracted from style '%s': %s\n", style, width)
		} else {
			fmt.Printf("Error extracting width from style: %s\n", style)
		}
	}
}

func TestReplace2(t *testing.T) {
	s := `首发 x 唐忠汉｜453㎡艺术大宅，优雅的多样诠释！ 

豪宅
453㎡优雅大宅
奔赴家的理想，
和梦境交叠，
如何在现实中完成梦想，
勾勒出内心的情怀？
从古至今，
家的形式不断变化，
时尚元素让人眼花缭乱。
经典美学被提炼，
结合当下的审美，
经得起岁月的沉淀，
这才是生活的幸福永恒。
大家好，我是《顶豪》的主理人豪叔，专为设计师分享全球顶级豪宅我们致力于激发设计师的创意火花，提供丰富多样的设计思路，推动中国设计走出国门。
今天来到一处453㎡的豪宅中，感受园林风光和惬意生活吧。
入院赏“园”
南山诗意 尽收眼底
居家生活和繁华城市之间的理想过渡，莫过于一处院落，这里充满想象力。设计师打造一处园林，绿植和水景契合自然。随心漫步，忘却时光，内心只剩下欣喜自在。


一处开放式茶室作为独处和会客的场所，茶香气和鸟语花香结合，彰显中式经典文化内涵。水景雅致，茶香怡人，充满诗情画意。




园林式打造手法，让院落超越功能本身，艺术气息烘托建筑的美感。每一个角落，都有独特的韵味，平时养花种草，也乐在其中。田园生活，自然归心，质朴阳光。


舒适的院落生活，精致的水景，勾勒出经典的中式意境。空间内则彰显奢华，深色的内敛和大气，彰显豪宅魅力。


跳出中式框架
暗黑风彰显优雅情怀
从自然景致的院落进入室内，脑海中还回味着中式情怀，然而设计师带来惊喜。巧妙的装饰，打造现代人居美景，奢华内敛的气场，线条和物件在空间勾勒。


黑色背景让光有了深邃的背景，打破框架束缚。黑色装饰在光的反馈下，展示出优雅一面。地下室的明亮感觉，没有任何压抑感，尽情体验功能趣味。
灯光设计是空间氛围的催化剂，不同位置的结构布局，搭配灯光设计，优雅从容。










黑色的典雅让不同灯光有了归宿，也衬托出这份鲜艳和雅致。明暗交错的空间内，几何美学彰显魅力，材质的肌理和光泽不再隐藏。






幽暗的环境下，光点缀一抹光亮，那绿植的新鲜绿意，直达内心深处。隐藏于心的情怀被x现实感染，迫不及待的奔赴赋予灵魂力量，营造人居惊喜。


木质肌理在灰色中不断勾勒，自然的笔法超越人工的精彩。绿植总会带来惊喜，仿佛一处独立空间，向周围晕染。艺术在隐约中释放精彩。


健身房运动气息浓郁，椭圆形天花板灯带设计，和地面设计契合。圈内是一个单车，随时爆发力量。木质的温馨和现代时尚元素融合。


于内映心意方圆明暗 温暖渐进
一处工作室，满足居者的爱好，没有多余的装饰，都是居者喜爱的物件和作品。桌子后边堆满了艺术品，这些都是艺术的沉淀。匠心精神不断熏陶，把美好传递。这些个性和艺术，彰显情怀。




来到地上一层空间，弧形沙发布局轻柔优雅，深色木质地板承载阳光的温暖。深色单人椅柔软舒适，也作为空间装饰。松弛的布局带来时光感觉，茶几简约精致，生活场景十分接地气。豪宅的人居保持高压姿态，也保持对生活和自然的尊重，沉淀美好回忆。








圆形灯具勾勒艳阳高照，光泽地面反射的倒影如水中明月。方圆意境彰显中式情怀。内心沉浸其中，感受情绪的变化，悠然自得。


卧室设计强调功能性，步入式衣帽间带来时尚感觉。柜面配色深邃，保持私密属性。递进式美学释放休息氛围，空间表达自我，享受生活。








室内车库设计，充满机械力量感觉，速度和激情得到展现。玻璃元素带来整洁效果。




拥有自然美景，园林随时漫步人生。内部的优雅超越想象，古典和现代结合，文化气息弥漫。生活也保持纯粹的状态，个人爱好添加其中。从手工室到车库，设计师考虑面面俱到。独特的艺术风格打造现代超前艺术场景，大胆勾勒生活唯美状态。
`

	fmt.Println(ReduceEmptyLines(s, 3))
}

// TrimString 裁剪字符串
func TrimString(s string, length int) string {
	var count int
	var result []rune
	for _, r := range s {
		if count >= length {
			break
		}
		if unicode.Is(unicode.Han, r) {
			count += 2
		} else {
			count++
		}
		result = append(result, r)
	}
	return string(result)
}

func TestCut(t *testing.T) {
	str := "你好我ello"
	newStr := TrimString(str, 5)
	fmt.Println(newStr)
}
