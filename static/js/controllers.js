(function(){
  var builderControllers = angular.module('builderControllers', []);

  builderControllers.controller('BuildModalCtrl', [ '$scope', '$http', '$modalInstance', '$location',
    function($scope, $http, $modalInstance, $location){

      $scope.submit = function() {
        build = {
          "clone_url":$scope.cloneURL,
          "commit_id":$scope.commitID
        }

        $http.post('/builds', build).success(function(data){
          $location.path("/builds/"+data.id);
          $modalInstance.close();
        });
      };

      $scope.cancel = function() {
        $modalInstance.close();
      };
    }]);

  builderControllers.controller('BuildListCtrl', [ '$scope', '$http', '$interval',
    function($scope, $http, $interval){
      $scope.builds = [];

      $scope.loadData = function(){
        $http.get('/builds').success(function(data){
          $scope.builds = data;
        });
      };

      $scope.loadData();

      $scope.loadLoop = $interval(function() {
        $scope.loadData();
      }, 3000);

      $scope.$on("$destroy", function(){
        $interval.cancel($scope.loadLoop);
      });
    }]);

  builderControllers.controller('BuildDetailCtrl', ['$scope', '$routeParams', '$http',
    function($scope, $routeParams, $http) {
      $scope.buildId = $routeParams.buildId;
      $scope.buildData = '';
      $scope.pushData = '';

      $scope.addStreamingData = function (msg) {
        $scope.$apply(function () { $scope.buildData = $scope.buildData + msg.data + "\n"; });
      };
      $scope.addStaticData = function (msg) {
        $scope.$apply(function () { $scope.pushData = msg.data; });
      };

      $scope.listen = function () {
        $scope.buildFeed = new EventSource('/builds/' + $scope.buildId + '/streams/build');
        $scope.buildFeed.addEventListener('message', $scope.addStreamingData, false);
        $scope.buildFeed.addEventListener('error', function(e) {
          $scope.buildFeed.close();
          $http.get('/builds/' + $scope.buildId + '/logs/build').success(function(data){
            $scope.buildData = data;
          });
        }, false);

        $scope.pushFeed = new EventSource('/builds/' + $scope.buildId + '/streams/push');
        $scope.pushFeed.addEventListener('message', $scope.addStaticData, false);
        $scope.pushFeed.addEventListener('error', function(e) {
          $scope.pushFeed.close();
          $http.get('/builds/' + $scope.buildId + '/logs/push').success(function(data){
            $scope.pushData = data;
          });
        }, false);
      };

      $scope.listen();

      $scope.$on("$destroy", function(){
        $scope.buildFeed.close();
        $scope.pushFeed.close();
      });
    }]);

})();
