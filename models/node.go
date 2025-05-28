package models

import (
	"gorm.io/gorm/clause"
)

type Node struct {
	ID              int `gorm:"primaryKey"`
	Link            string
	Name            string `gorm:"uniqueIndex:idx_name_id"`
	DialerProxyName string
	CreateDate      string
	Source          string `gorm:"default:'manual'"`
}

// Add 添加节点
func (node *Node) Add() error {
	return DB.Create(node).Error
}

// 更新节点
func (node *Node) Update() error {
	return DB.Model(node).Select("Name", "Link", "DialerProxyName").Updates(node).Error
}

// 查找节点是否重复
func (node *Node) Find() error {
	return DB.Where("link = ? or name = ?", node.Link, node.Name).First(node).Error
}

// 节点列表
func (node *Node) List() ([]Node, error) {
	var nodes []Node
	err := DB.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// 删除节点
func (node *Node) Del() error {
	return DB.Delete(node).Error
}

// UpsertNode 插入或更新节点
func (node *Node) UpsertNode() error {
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"link", "create_date", "source"}),
	}).Create(node).Error
}

// DeleteAutoSubscriptionNodes 删除订阅节点
func DeleteAutoSubscriptionNodes(subName string) error {
	return DB.Where("source = ?", "sublinkE").Where("name like ?", subName+"%").Delete(&Node{}).Error
}
