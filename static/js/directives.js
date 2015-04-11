(function(){
  var app = angular.module('builderDirectives', []);

  app.directive("navbar", function() {
    return {
      restrict: "E",
      templateUrl: "partials/nav.html",
      controller: function($scope, $modal) {
        $scope.newBuild = function() {
          var modalInstance = $modal.open({
            templateUrl: 'partials/build-new.html',
            controller: 'BuildModalCtrl',
          });
        };
      },
      controllerAs: "navbar"
    };
  });

})();
