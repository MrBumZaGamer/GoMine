package worlds

import (
	"gomine/interfaces"
	"os"
	"errors"
)

type LevelManager struct {
	server interfaces.IServer
	levels map[string]interfaces.ILevel
}

func NewLevelManager(server interfaces.IServer) *LevelManager {
	return &LevelManager{server, make(map[string]interfaces.ILevel)}
}

/**
 * Returns all loaded levels in the manager.
 */
func (manager *LevelManager) GetLoadedLevels() map[string]interfaces.ILevel {
	return manager.levels
}

/**
 * Returns whether a level is loaded or not.
 */
func (manager *LevelManager) IsLevelLoaded(levelName string) bool {
	var _, ok = manager.levels[levelName]
	return ok
}

/**
 * Returns whether a level is generated or not. (Includes loaded levels)
 */
func (manager *LevelManager) IsLevelGenerated(levelName string) bool {
	if manager.IsLevelLoaded(levelName) {
		return true
	}
	var path = manager.server.GetServerPath() + "worlds/" + levelName
	var _, err = os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

/**
 * Loads a generated world. Returns true if the level was loaded successfully.
 */
func (manager *LevelManager) LoadLevel(levelName string) bool {
	if !manager.IsLevelGenerated(levelName) {
		// manager.GenerateLevel(level) We need file writing for this. TODO.
	}
	if manager.IsLevelLoaded(levelName) {
		return false
	}
	manager.levels[levelName] = NewLevel(levelName, manager.server)
	return true
}

/**
 * Returns the default level and loads/generates it if needed.
 */
func (manager *LevelManager) GetDefaultLevel() interfaces.ILevel {
	if !manager.IsLevelLoaded(manager.server.GetConfiguration().DefaultLevel) {
		manager.LoadLevel(manager.server.GetConfiguration().DefaultLevel)
	}
	var level, _ = manager.GetLevelByName(manager.server.GetConfiguration().DefaultLevel)
	return level
}

/**
 * Returns a level by its name. Returns an error if the level is not loaded.
 */
func (manager *LevelManager) GetLevelByName(name string) (interfaces.ILevel, error) {
	var level interfaces.ILevel
	if !manager.IsLevelGenerated(name) {
		return level, errors.New("level with given name is not generated")
	}
	if !manager.IsLevelLoaded(name) {
		return level, errors.New("level with given name is not loaded")
	}

	return manager.levels[name], nil
}

/**
 * Ticks all levels in the manager.
 */
func (manager *LevelManager) Tick() {
	for _, level := range manager.levels {
		level.TickLevel()
	}
}

func (manager *LevelManager) Close() {
	for _, level := range manager.levels {
		for _, dimension := range level.GetDimensions() {
			dimension.Close()
		}
	}
}