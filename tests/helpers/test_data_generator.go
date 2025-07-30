package testhelpers

import (
	"fmt"
	"math/rand"
	"time"
)

// TestDataGenerator 测试数据生成器
type TestDataGenerator struct {
	rnd *rand.Rand
}

// NewTestDataGenerator 创建新的测试数据生成器
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateRandomString 生成指定长度的随机字符串
func (g *TestDataGenerator) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[g.rnd.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomUsername 生成随机用户名
func (g *TestDataGenerator) GenerateRandomUsername() string {
	timestamp := time.Now().UnixNano()
	randomSuffix := g.GenerateRandomString(6)
	return fmt.Sprintf("testuser_%d_%s", timestamp, randomSuffix)
}

// GenerateRandomEmail 生成随机邮箱
func (g *TestDataGenerator) GenerateRandomEmail() string {
	timestamp := time.Now().UnixNano()
	randomSuffix := g.GenerateRandomString(6)
	return fmt.Sprintf("test_%d_%s@example.com", timestamp, randomSuffix)
}

// GenerateRandomFullName 生成随机全名
func (g *TestDataGenerator) GenerateRandomFullName() string {
	firstNames := []string{"张", "李", "王", "刘", "陈", "杨", "赵", "黄", "周", "吴"}
	lastNames := []string{"伟", "芳", "娜", "敏", "静", "丽", "强", "磊", "军", "洋"}
	
	firstName := firstNames[g.rnd.Intn(len(firstNames))]
	lastName := lastNames[g.rnd.Intn(len(lastNames))]
	
	return firstName + lastName
}

// GenerateRandomDiaryTitle 生成随机日记标题
func (g *TestDataGenerator) GenerateRandomDiaryTitle() string {
	titles := []string{
		"今天的美好时光",
		"工作中的小确幸",
		"周末的悠闲时光",
		"学习新技能的感悟",
		"与朋友的聚会",
		"读书笔记",
		"旅行见闻",
		"美食体验",
		"运动健身记录",
		"生活感悟",
	}
	
	baseTitle := titles[g.rnd.Intn(len(titles))]
	timestamp := time.Now().Format("0102_1504")
	
	return fmt.Sprintf("%s_%s", baseTitle, timestamp)
}

// GenerateRandomDiaryContent 生成随机日记内容
func (g *TestDataGenerator) GenerateRandomDiaryContent() string {
	contents := []string{
		"今天过得很充实，学到了很多新东西。",
		"天气很好，心情也很不错。",
		"工作虽然忙碌，但很有成就感。",
		"和家人一起度过了愉快的时光。",
		"读了一本好书，收获颇丰。",
		"尝试了新的菜谱，味道还不错。",
		"运动后感觉身体很舒服。",
		"解决了一个技术难题，很有成就感。",
		"欣赏了美丽的日落，心情很平静。",
		"计划明天要做的事情，充满期待。",
	}
	
	content := contents[g.rnd.Intn(len(contents))]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	return fmt.Sprintf("%s\n\n记录时间: %s", content, timestamp)
}

// GenerateRandomPostTitle 生成随机文章标题
func (g *TestDataGenerator) GenerateRandomPostTitle() string {
	titles := []string{
		"技术分享",
		"生活感悟",
		"学习笔记",
		"项目总结",
		"经验分享",
		"问题解决方案",
		"工具推荐",
		"读书心得",
		"旅行游记",
		"美食推荐",
	}
	
	baseTitle := titles[g.rnd.Intn(len(titles))]
	timestamp := time.Now().Format("0102_1504")
	randomSuffix := g.GenerateRandomString(4)
	
	return fmt.Sprintf("%s_%s_%s", baseTitle, timestamp, randomSuffix)
}

// 全局测试数据生成器实例
var globalTestDataGenerator = NewTestDataGenerator()

// 便捷函数
func GenerateRandomUsername() string {
	return globalTestDataGenerator.GenerateRandomUsername()
}

func GenerateRandomEmail() string {
	return globalTestDataGenerator.GenerateRandomEmail()
}

func GenerateRandomFullName() string {
	return globalTestDataGenerator.GenerateRandomFullName()
}

func GenerateRandomDiaryTitle() string {
	return globalTestDataGenerator.GenerateRandomDiaryTitle()
}

func GenerateRandomDiaryContent() string {
	return globalTestDataGenerator.GenerateRandomDiaryContent()
}

func GenerateRandomPostTitle() string {
	return globalTestDataGenerator.GenerateRandomPostTitle()
}