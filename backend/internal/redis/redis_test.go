package redis_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	crdt "ugboss/crdt/internal/redis"
)

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Suite")
}

var _ = Describe("redis", func() {
	var  (
		cli *redis.Client
		ctx = context.TODO()
		list = "list"
		deleted = list + "-deleted"
	)

	BeforeEach(func() {
		cli = redis.NewClient(&redis.Options{
			Addr: ":6379",
			Password: "",
			DB: 0,
		})
		Expect(cli.FlushDB(ctx).Err()).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		Expect(cli.Close()).NotTo(HaveOccurred())
	})

	Describe("range from to", func() {
		It("should return all", func() {
			cli.RPush(ctx, list, 1, 2, 3, 4, 5)

			res, err := crdt.RangeFromTo(ctx, cli, list, 0, 4)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal([]string{"1", "2", "3", "4", "5"}))
		})
	})

	Describe("range from", func() {
		It("should return all", func() {
			cli.RPush(ctx, list, 1, 2, 3, 4, 5)
			
			res, err := crdt.RangeFrom(ctx, cli, list, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal([]string{"1", "2", "3", "4", "5"}))
		})

		It("should return 5 to 8", func() {
			cli.RPush(ctx, list, 4, 5, 6, 7, 8)
			cli.Set(ctx, deleted, 3, 0)
			
			res, err := crdt.RangeFrom(ctx, cli, list, 4)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal([]string{"5", "6", "7", "8"}))	
		})

		It("should occure error", func() {
			cli.RPush(ctx, list, 4, 5, 6, 7, 8)
			cli.Set(ctx, deleted, 3, 0)
			
			_, err := crdt.RangeFrom(ctx, cli, list, 1)
			Expect(err).To(HaveOccurred())
		})
	})
})
