(function(){
  var builderControllers = angular.module('builderControllers', []);

  builderControllers.controller('BuildModalCtrl', [ '$scope', '$http', '$modalInstance',
    function($scope, $http, $modalInstance){

      $scope.submit = function() {
        build = {
          "repository_name":$scope.repositoryName,
          "clone_url":$scope.cloneURL,
          "commit_id":$scope.commitID,
          "git_ref":$scope.gitRef
        }
        console.log(build)
        $http.post('http://localhost:3000/builds', build).success(function(data){
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
        $http.get('http://localhost:3000/builds').success(function(data){
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

  builderControllers.controller('BuildDetailCtrl', ['$scope', '$routeParams',
    function($scope, $routeParams) {
      $scope.buildId = $routeParams.buildId;
      $scope.buildData = '';

      $scope.addData = function (msg) {
        $scope.$apply(function () { $scope.buildData = $scope.buildData + msg.data + "\n"; });
      };

      $scope.listen = function () {
        $scope.buildFeed = new EventSource('/builds/' + $scope.buildId);
        $scope.buildFeed.addEventListener('message', $scope.addData, false);
      };

      $scope.listen();

      $scope.$on("$destroy", function(){
        $scope.buildFeed.close();
      });
    }]);

})();
