package main

import (
	"context"

	"github.com/segmentio/cycle"
	"github.com/segmentio/events"
)

// envLogs is an Environment decorator which adds logging to every
// calls to methods of the base environment.
type envLogs struct {
	base cycle.Environment
}

func (env envLogs) LookupClusterID(ctx context.Context, name string) (cycle.ClusterID, error) {
	clusterID, err := env.base.LookupClusterID(ctx, name)
	if err != nil {
		events.Log("error looking up %{cluster_name}s cluster - %{error}v", name, err)
	} else {
		events.Debug("%{cluster_id}s - found %{cluster_name}s cluster", clusterID, name)
	}
	return clusterID, err
}

func (env envLogs) DescribeCluster(ctx context.Context, id cycle.ClusterID) (cycle.Cluster, error) {
	cluster, err := env.base.DescribeCluster(ctx, id)
	if err != nil {
		events.Log("%{cluster_id}s - error describing cluster - %{error}v", id, err)
	} else {
		outdated := 0

		for _, instance := range cluster.Instances {
			if instance.Config != cluster.Config {
				outdated++
			}
		}

		events.Debug("%{cluster_id}s - found configuration %{config_id}s and %{outdated_instance_count}d/%{instance_count}d outdated instances (min size: %d, max size: %d)",
			id, cluster.Config, outdated, len(cluster.Instances), cluster.MinSize, cluster.MaxSize)
	}
	return cluster, err
}

func (env envLogs) StartInstances(ctx context.Context, cluster cycle.ClusterID, count int) error {
	events.Debug("%{cluster_id}s - starting %{instance_count}d new instances", cluster, count)
	err := env.base.StartInstances(ctx, cluster, count)
	if err != nil {
		events.Log("error starting instances - %{error}v", err)
	}
	return err
}

func (env envLogs) DrainInstances(ctx context.Context, instances ...cycle.InstanceID) error {
	for _, instance := range instances {
		events.Debug("%{instance_id}s - draining", instance)
	}
	err := env.base.DrainInstances(ctx, instances...)
	if err != nil {
		events.Log("error draining instances - %{error}v", err)
	}
	return err
}

func (env envLogs) TerminateInstances(ctx context.Context, instances ...cycle.InstanceID) error {
	for _, instance := range instances {
		events.Debug("%{instance_id}s - terminating", instance)
	}
	err := env.base.TerminateInstances(ctx, instances...)
	if err != nil {
		events.Log("error terminating instances - %{error}v", err)
	}
	return nil
}

func (env envLogs) WaitInstances(ctx context.Context, state cycle.InstanceState, instances ...cycle.InstanceID) error {
	for _, instance := range instances {
		events.Debug("%{instance_id}s - waiting to be %{waiting_state}s", instance, state)
	}
	err := env.base.WaitInstances(ctx, state, instances...)
	if err != nil {
		events.Log("error waiting for instance to be %{waiting_state}s - %{error}v", state, err)
	}
	return err
}
