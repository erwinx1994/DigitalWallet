package implementation

import (
	"api_gateway/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redis_manager struct {
	deposit_requests_queue              *redis.Client
	deposit_responses_queue             *redis.Client
	withdrawal_requests_queue           *redis.Client
	withdrawal_responses_queue          *redis.Client
	transfer_requests_queue             *redis.Client
	transfer_responses_queue            *redis.Client
	balance_requests_queue              *redis.Client
	balance_responses_queue             *redis.Client
	transaction_history_requests_queue  *redis.Client
	transaction_history_responses_queue *redis.Client
}

func create_redis_manager(config *config.Config, background_context context.Context) (*redis_manager, error) {
	redis_manager := &redis_manager{}

	{
		// Prepare deposit requests queue
		deposit_requests_queue := redis.NewClient(config.DepositsService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.DepositsService.RequestsQueue.Timeout)*time.Second)
		_, err := deposit_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.deposit_requests_queue = deposit_requests_queue
	}

	{
		// Prepare deposit responses queue
		deposit_responses_queue := redis.NewClient(config.DepositsService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.DepositsService.ResponsesQueue.Timeout)*time.Second)
		_, err := deposit_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.deposit_responses_queue = deposit_responses_queue
	}

	{
		// Prepare withdrawal requests queue
		withdrawal_requests_queue := redis.NewClient(config.WithdrawalService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.WithdrawalService.RequestsQueue.Timeout)*time.Second)
		_, err := withdrawal_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.withdrawal_requests_queue = withdrawal_requests_queue
	}

	{
		// Prepare withdrawal responses queue
		withdrawal_responses_queue := redis.NewClient(config.WithdrawalService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.WithdrawalService.ResponsesQueue.Timeout)*time.Second)
		_, err := withdrawal_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.withdrawal_responses_queue = withdrawal_responses_queue
	}

	{
		// Prepare transfer requests queue
		transfer_requests_queue := redis.NewClient(config.TransferService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.TransferService.RequestsQueue.Timeout)*time.Second)
		_, err := transfer_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.transfer_requests_queue = transfer_requests_queue
	}

	{
		// Prepare transfer responses queue
		transfer_responses_queue := redis.NewClient(config.TransferService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.TransferService.ResponsesQueue.Timeout)*time.Second)
		_, err := transfer_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.transfer_responses_queue = transfer_responses_queue
	}

	{
		// Prepare balance requests queue
		balance_requests_queue := redis.NewClient(config.BalanceService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.BalanceService.RequestsQueue.Timeout)*time.Second)
		_, err := balance_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.balance_requests_queue = balance_requests_queue
	}

	{
		// Prepare balance responses queue
		balance_responses_queue := redis.NewClient(config.BalanceService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.BalanceService.ResponsesQueue.Timeout)*time.Second)
		_, err := balance_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.balance_responses_queue = balance_responses_queue
	}

	{
		// Prepare transaction history requests queue
		transaction_history_requests_queue := redis.NewClient(config.TransactionHistoryService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.TransactionHistoryService.RequestsQueue.Timeout)*time.Second)
		_, err := transaction_history_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.transaction_history_requests_queue = transaction_history_requests_queue
	}

	{
		// Prepare transaction history responses queue
		transaction_history_responses_queue := redis.NewClient(config.TransactionHistoryService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.TransactionHistoryService.ResponsesQueue.Timeout)*time.Second)
		_, err := transaction_history_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()

		redis_manager.transaction_history_responses_queue = transaction_history_responses_queue
	}

	return redis_manager, nil
}
